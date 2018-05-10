package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	// "strings"
	"encoding/gob"
	"github.com/gorilla/sessions"
	"net/rpc"
	"os"
	// "reflect"
	"time"
)

var (
	// key must be 16, 24 or 32 bytes long (AES-128, AES-192 or AES-256)
	key   = []byte("super-secret-key")
	store = sessions.NewCookieStore(key)

	// back end server address
	serverAddress string = "localhost"
)

type PageVariables struct {
	Date string
	Time string
}

type HomePageArgs struct {
	User             UserProfile
	MsgFromFollowing []Message
}

type GetProfilePageArgs struct {
	User, CurrentUser string
}

type GetProfilePageReply struct {
	isFollowing bool
	User        UserProfile
}

type UserProfile struct {
	UserEmail string
	UserId    int
	UserName  string
	UserBio   string
	// Following []int
	// Follower  []int
	// PostMsg   []Message
	FollowingNum int
	FollowerNum  int
}

type UserCredential struct {
	UserEmail string
	Password  string
}

type Message struct {
	SenderEmail string
	SenderName  string
	// MsgTimeStamp Time
	MsgContent string
}

type LoginArgs struct {
	UserLoginEmail, UserLoginPassword string
}

type LoginReply struct {
	UserLoginStatus  bool
	UserLoginProfile UserProfile
}

type SignUpArgs struct {
	UserSignUpEmail, UserSignUpName, UserSignUpPassword string
}

type SignUpReply struct {
	UserSignUpStatus  bool
	UserSignUpProfile UserProfile
}

// type SendMessageArgs struct {
// 	SendMsgEmail, SendMsgName string
// 	Msg                       Message
// }

type SendMessageReply struct {
	SendMsgStatus  bool
	SendMsgProfile UserProfile
}

func init() {
	gob.Register(&UserProfile{})
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: ", os.Args[0], "server")
		os.Exit(1)
	}
	serverAddress = os.Args[1]

	//
	http.Handle("/resources/css/", http.StripPrefix("/resources/css/", http.FileServer(http.Dir("resources/css"))))

	//root
	http.HandleFunc("/", indexHandler)
	// log.Fatal(http.ListenAndServe(":8080", nil))

	//Login
	http.HandleFunc("/login", loginHandler)

	// signUpRedirect
	http.HandleFunc("/signUpRedirect", signUpRedirect)

	//SignUp
	http.HandleFunc("/signup", signUp)

	//sendMessage
	http.HandleFunc("/sendMessage", sendMessage)

	//logout
	http.HandleFunc("/logout", logout)

	//searchUsers
	http.HandleFunc("/searchUsers", searchUsersHandler)

	// getUserProfile
	http.HandleFunc("/getUserProfile", getUserProfileHandler)

	err := http.ListenAndServe(":8080", nil) // setting listening port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {

	now := time.Now()              // find the time right now
	HomePageVars := PageVariables{ //store the date and time in a struct
		Date: now.Format("02-01-2006"),
		Time: now.Format("15:04:05"),
	}

	t, err := template.ParseFiles("index.html") //parse the html file homepage.html
	if err != nil {                             // if there is an error
		log.Print("Index template parsing error: ", err) // log it
	}
	err = t.Execute(w, HomePageVars) //execute the template and pass it the HomePageVars struct to fill in the gaps
	if err != nil {                  // if there is an error
		log.Print("Index template executing error: ", err) //log it
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Session: User is Logging in.")
	fmt.Println("method:", r.Method) //get request method
	if r.Method == "GET" {
		t, _ := template.ParseFiles("login.gtpl")
		t.Execute(w, nil)
	} else {
		r.ParseForm()
		// fmt.Printf("%+v\n", r.Form)
		// for key, values := range r.Form { // range over map
		// 	for _, value := range values { // range over []string
		// 		fmt.Println(key, value)
		// 	}
		// }
		// logic part of log in
		userEmail := r.FormValue("useremail")
		userPswd := r.FormValue("userpswd")
		// fmt.Println("userName : ", userName)
		// fmt.Println("password : ", password)

		session, err := store.Get(r, "cookie-name")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Authentication goes here
		client, err := rpc.DialHTTP("tcp", serverAddress+":1234")
		if err != nil {
			log.Fatal("dialing:", err)
		}
		// Synchronous call
		// RPC call for validation
		loginArgs := LoginArgs{userEmail, userPswd}
		var loginReply LoginReply
		err = client.Call("UserProfile.UserLoginHandler", loginArgs, &loginReply)
		if err != nil {
			log.Fatal("User Login error:", err)
		}
		fmt.Printf("User: %s, LoginStatus: %t\n", loginArgs.UserLoginEmail, loginReply.UserLoginStatus)

		if loginReply.UserLoginStatus == true {
			// Set user as authenticated
			session.Values["authenticated"] = true
			// Retrieve our struct and type-assert it
			// val := session.Values["user"]
			// var usr = &UserProfile{}
			// if usr, ok := val.(*UserProfile); !ok {
			// 	// Handle the case that it's not an expected type
			// 	fmt.Printf("Error in signUp\n")
			// }
			// fmt.Println(usr)
			// Now we can use our User object
			session.Values["currentUser"] = loginReply.UserLoginProfile

			session.Save(r, w)

			// tmpl, err := template.ParseFiles("templates/homepage.html") //parse the html file homepage.html
			// if err != nil {                                             // if there is an error
			// 	log.Print("template parsing error: ", err) // log it
			// }

			// tmpl.Execute(w, loginReply.UserLoginProfile)
			getHomePage(w, r, &loginReply.UserLoginProfile)
		} else {
			tmpl := template.Must(template.ParseFiles("templates/signUp.html"))
			tmpl.Execute(w, nil)
			// fmt.Fprintf(w, "Sorry, %s. Sign Up and Join us today.\n", userName[0])
		}

	}
}

func getHomePage(w http.ResponseWriter, r *http.Request, usr *UserProfile) {
	fmt.Println("Session: Getting HomePage.")
	if r.Method == "GET" {
		t, _ := template.ParseFiles("login.gtpl")
		t.Execute(w, nil)
	} else {
		tmpl, err := template.ParseFiles("templates/homepage.html") //parse the html file homepage.html
		if err != nil {                                             // if there is an error
			log.Print("template parsing error: ", err) // log it
		}

		// Start getting Msg
		client, err := rpc.DialHTTP("tcp", serverAddress+":1234")
		if err != nil {
			log.Fatal("dialing:", err)
		}
		// Synchronous call
		// RPC call

		var getMsgReply []Message
		err = client.Call("UserProfile.GetUserMsgFromFollowingHandler", usr, &getMsgReply)
		if err != nil {
			log.Fatal("User Login error:", err)
		}

		if getMsgReply == nil {
			fmt.Println("GetUserMsgFromFollowing: Error, Redirect User to Sign Up.")
			signUpRedirect(w, r)
		} else {
			hpArgs := HomePageArgs{
				User:             *usr,
				MsgFromFollowing: getMsgReply,
			}
			tmpl.Execute(w, hpArgs)
		}

	}
}

func signUpRedirect(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Session: Redirect User to Sign Up.")
	if r.Method == "GET" {
		t, _ := template.ParseFiles("login.gtpl")
		t.Execute(w, nil)
	} else {
		tmpl := template.Must(template.ParseFiles("templates/signUp.html"))
		tmpl.Execute(w, nil)
	}
}

func signUp(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Session: User is Signing Up.")
	fmt.Println("method:", r.Method) //get request method
	if r.Method == "GET" {
		t, _ := template.ParseFiles("login.gtpl")
		t.Execute(w, nil)
	} else {
		r.ParseForm()
		// logic part of log in
		userEmail := r.FormValue("useremail")
		userName := r.FormValue("username")
		userPswd := r.FormValue("userpswd")

		// SignUp RPCs
		client, err := rpc.DialHTTP("tcp", serverAddress+":1234")
		if err != nil {
			log.Fatal("dialing:", err)
		}
		// Synchronous call
		signUpArgs := SignUpArgs{userEmail, userName, userPswd}
		var signUpReply SignUpReply
		err = client.Call("UserProfile.UserSignUpHandler", signUpArgs, &signUpReply)
		if err != nil {
			log.Fatal("User SignUp error:", err)
		}
		fmt.Printf("User: %s, SignUpStatus: %t\n", signUpArgs.UserSignUpName, signUpReply.UserSignUpStatus)

		if signUpReply.UserSignUpStatus == false {
			//userName already used!
			tmpl := template.Must(template.ParseFiles("templates/signUp.html"))
			tmpl.Execute(w, nil)
		} else {
			//signUp success
			//Redirect to index
			indexHandler(w, r)
			// fmt.Fprintf(w, "Sorry, %s. Sign Up and Join us today.\n", userName[0])
		}

	}
}

func sendMessage(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")

	// Check if user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	fmt.Println("Session: User is Sending Message.")
	fmt.Println("method:", r.Method) //get request method
	if r.Method == "GET" {
		t, _ := template.ParseFiles("login.gtpl")
		t.Execute(w, nil)
	} else {
		r.ParseForm()

		if userPfl, ok := session.Values["currentUser"].(*UserProfile); !ok || userPfl == nil {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		userPfl := session.Values["currentUser"].(*UserProfile)

		newMsg := Message{
			SenderEmail: userPfl.UserEmail,
			SenderName:  userPfl.UserName,
			MsgContent:  r.FormValue("message"),
		}
		// fmt.Fprintf(w, "Message <div> <p>%s</p> </div> sent.\n", msg)
		fmt.Println(newMsg)

		// sendMessage RPCs
		client, err := rpc.DialHTTP("tcp", serverAddress+":1234")
		if err != nil {
			log.Fatal("dialing:", err)
		}
		// Synchronous call
		msgArgs := newMsg
		var msgReply SendMessageReply
		err = client.Call("UserProfile.SendMessageHandler", msgArgs, &msgReply)
		if err != nil {
			log.Fatal("Send Message error:", err)
		}
		fmt.Printf("User: %s, sendMessageStatus: %t\n", msgArgs.SenderName, msgReply.SendMsgStatus)

		//re-render the page
		if msgReply.SendMsgStatus == true {
			// tmpl := template.Must(template.ParseFiles("templates/homepage.html"))
			// tmpl.Execute(w, msgReply.SendMsgProfile)
			getHomePage(w, r, userPfl)
		} else {
			fmt.Println("SendMessage: userName does not exist, go to sign up.")
			logout(w, r)
		}

		// fmt.Fprintf(w, "Sorry, %s. Sign Up and Join us today.\n", userName[0])
	}
}

func searchUsersHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")

	// Check if user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	fmt.Println("Session: User is Searching Users.")
	fmt.Println("method:", r.Method) //get request method
	if r.Method == "GET" {
		t, _ := template.ParseFiles("login.gtpl")
		t.Execute(w, nil)
	} else {
		r.ParseForm()

		if userPfl, ok := session.Values["currentUser"].(*UserProfile); !ok || userPfl == nil {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		// userPfl := session.Values["currentUser"].(*UserProfile)

		searchKeyword := r.FormValue("searchkeyword")

		// searchResults RPCs
		client, err := rpc.DialHTTP("tcp", serverAddress+":1234")
		if err != nil {
			log.Fatal("dialing:", err)
		}
		// Synchronous call
		searchArgs := searchKeyword
		var searchReply []UserProfile
		err = client.Call("UserProfile.GetSearchResultsHandler", searchArgs, &searchReply)
		if err != nil {
			log.Fatal("Search Users error:", err)
		}
		// fmt.Printf("User: %s, sendMessageStatus: %t\n", msgArgs.SenderName, msgReply.SendMsgStatus)

		//re-render the page
		// if searchReply != nil {
		tmpl, err := template.ParseFiles("templates/searchResults.html") //parse the html file homepage.html
		if err != nil {                                                  // if there is an error
			log.Print("template parsing error: ", err) // log it
		}

		tmpl.Execute(w, searchReply)
		// } else {
		// 	fmt.Println("SendMessage: userName does not exist, go to sign up.")
		// 	logout(w, r)
		// }

		// fmt.Fprintf(w, "Sorry, %s. Sign Up and Join us today.\n", userName[0])
	}
}

func getUserProfileHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")

	// Check if user is authenticated
	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	fmt.Println("Session: User is Getting Users Profile.")
	fmt.Println("method:", r.Method) //get request method
	if r.Method == "GET" {
		t, _ := template.ParseFiles("login.gtpl")
		t.Execute(w, nil)
	} else {
		r.ParseForm()

		if userPfl, ok := session.Values["currentUser"].(*UserProfile); !ok || userPfl == nil {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		currentUserProfl := session.Values["currentUser"].(*UserProfile)
		currentUserEmail := currentUserProfl.UserEmail

		userEmail := r.FormValue("userEmail")

		// getUserProfile RPCs
		client, err := rpc.DialHTTP("tcp", serverAddress+":1234")
		if err != nil {
			log.Fatal("dialing:", err)
		}
		// Synchronous call
		getUserProfileArgs := GetProfilePageArgs{
			userEmail,
			currentUserEmail,
		}
		var pfReply GetProfilePageReply
		err = client.Call("UserProfile.GetUserProfileHandler", getUserProfileArgs, &pfReply)
		if err != nil {
			log.Fatal("Get User Profile error:", err)
		}

		tmpl, err := template.ParseFiles("templates/profilePage.html") //parse the html file homepage.html
		if err != nil {                                                // if there is an error
			log.Print("template parsing error: ", err) // log it
		}

		fmt.Println(pfReply)

		if pfReply.isFollowing == true {
			fmt.Printf("%s is Following the user", currentUserProfl.UserName)
		} else {
			fmt.Printf("%s has not yet followed the user", currentUserProfl.UserName)
		}

		tmpl.Execute(w, pfReply)
	}
}

func logout(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Session: User is logging out.")
	session, _ := store.Get(r, "cookie-name")

	// Revoke users authentication
	session.Values["authenticated"] = false
	session.Save(r, w)

	//Redirect to index
	indexHandler(w, r)
}
