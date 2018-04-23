package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	// "strings"
	"github.com/gorilla/sessions"
	"net/rpc"
	"os"
	"reflect"
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

type User struct {
	UserId    int
	UserName  string
	Password  string
	Following []int
}

type Message struct {
	SenderId int
	Content  string
}

var userList = []User{
	{UserId: 1,
		UserName:  "Ned Stark",
		Password:  "qwerty",
		Following: []int{2}},
	{UserId: 2,
		UserName:  "Robert Baratheon",
		Password:  "qwerty",
		Following: []int{1}},
	{UserId: 3,
		UserName:  "Jaime Lannister",
		Password:  "qwerty",
		Following: []int{1, 2, 7}},
	{UserId: 4,
		UserName:  "Jon Snow",
		Password:  "qwerty",
		Following: []int{1, 6, 7}},
	{UserId: 5,
		UserName:  "Tyrion Lannister",
		Password:  "qwerty",
		Following: []int{1}},
	{UserId: 6,
		UserName:  "Daenerys Targaryen",
		Password:  "qwerty",
		Following: []int{4, 7}},
	{UserId: 7,
		UserName:  "Cersei Lannister",
		Password:  "qwerty",
		Following: []int{3, 6}},
}

type Args struct {
	A, B int
}

type Quotient struct {
	Quo, Rem int
}

type LoginArgs struct {
	UserLoginName, UserLoginPassword string
}

type LoginReply struct {
	UserLoginStatus  bool
	UserLoginProfile User
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
	http.HandleFunc("/", Index)
	// log.Fatal(http.ListenAndServe(":8080", nil))

	//Login
	http.HandleFunc("/login", login)

	//SignUp
	http.HandleFunc("/signup", signUp)

	//sendMessage
	http.HandleFunc("/sendMessage", sendMessage)

	//logout
	http.HandleFunc("/logout", logout)

	err := http.ListenAndServe(":8080", nil) // setting listening port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func Index(w http.ResponseWriter, r *http.Request) {

	now := time.Now()              // find the time right now
	HomePageVars := PageVariables{ //store the date and time in a struct
		Date: now.Format("02-01-2006"),
		Time: now.Format("15:04:05"),
	}

	t, err := template.ParseFiles("index.html") //parse the html file homepage.html
	if err != nil {                             // if there is an error
		log.Print("template parsing error: ", err) // log it
	}
	err = t.Execute(w, HomePageVars) //execute the template and pass it the HomePageVars struct to fill in the gaps
	if err != nil {                  // if there is an error
		log.Print("template executing error: ", err) //log it
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Session: User is Logging in.")
	fmt.Println("method:", r.Method) //get request method
	if r.Method == "GET" {
		t, _ := template.ParseFiles("login.gtpl")
		t.Execute(w, nil)
	} else {
		r.ParseForm()
		// logic part of log in
		fmt.Println("username:", r.Form["username"])
		fmt.Println("password:", r.Form["password"])
		userName := r.Form["username"]
		password := r.Form["passward"]
		fmt.Println(reflect.TypeOf(userName))
		fmt.Println(reflect.TypeOf(password))
		// fmt.Println("userName len: ", len(userName))
		// fmt.Println("password len: ", len(password))
		// fmt.Println("userName : ", userName)
		// fmt.Println("password : ", r.Form["passward"])

		session, _ := store.Get(r, "cookie-name")

		// Authentication goes here
		client, err := rpc.DialHTTP("tcp", serverAddress+":1234")
		if err != nil {
			log.Fatal("dialing:", err)
		}
		// Synchronous call
		// RPC call for validation
		tempPassword := "qwerty"
		loginArgs := LoginArgs{userName[0], tempPassword}
		fmt.Printf("I am on line 179\n")
		var loginReply LoginReply
		err = client.Call("User.UserLoginValidation", loginArgs, &loginReply)
		if err != nil {
			log.Fatal("User Login error:", err)
		}
		fmt.Printf("User: %s, LoginStatus: %t\n", loginArgs.UserLoginName, loginReply.UserLoginStatus)

		if loginReply.UserLoginStatus == true {
			// Set user as authenticated
			session.Values["authenticated"] = true
			session.Save(r, w)

			tmpl, err := template.ParseFiles("templates/homepage.html") //parse the html file homepage.html
			if err != nil {                                             // if there is an error
				log.Print("template parsing error: ", err) // log it
			}

			tmpl.Execute(w, loginReply.UserLoginProfile)
		} else {
			tmpl := template.Must(template.ParseFiles("templates/signUp.html"))
			tmpl.Execute(w, User{UserName: userName[0]})
			// fmt.Fprintf(w, "Sorry, %s. Sign Up and Join us today.\n", userName[0])
		}

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
		userName := r.Form["username"]
		// password := r.Form["passward"]

		//validate username
		loginStatus := true
		for _, v := range userList {
			if v.UserName == userName[0] {
				// Duplicate username
				loginStatus = false

				tmpl := template.Must(template.ParseFiles("templates/signUp.html"))
				tmpl.Execute(w, User{UserName: userName[0]})
			}

		}

		//login success
		newUserProfile := User{
			UserId:    len(userList) + 1,
			UserName:  userName[0],
			Password:  "qwerty",
			Following: []int{},
		}

		if loginStatus == true {
			session, _ := store.Get(r, "cookie-name")
			// Set user as authenticated
			session.Values["authenticated"] = true
			session.Save(r, w)

			tmpl := template.Must(template.ParseFiles("templates/homepage.html"))
			tmpl.Execute(w, newUserProfile)
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
		// logic part of log in
		msg := r.Form["message"]
		// password := r.Form["passward"]
		fmt.Fprintf(w, "Message <div> <p>%s</p> </div> sent.\n", msg)

	}
}

func logout(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookie-name")

	// Revoke users authentication
	session.Values["authenticated"] = false
	session.Save(r, w)
}
