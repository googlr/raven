package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	// "strings"
	"time"
)

type PageVariables struct {
	Date string
	Time string
}

type User struct {
  UserId int
	UserName  string
	Password  string
	Following []int
}

type Message struct {
	SenderId  int
	Content string
}


func main() {

  //create some users

  //
  userList := [] User {
    { UserId: 1, 
    UserName:"Ned Stark", 
    Password:"qwerty", 
    Following {2}},
  { UserId: 2, 
    UserName:"Robert Baratheon", 
    Password:"qwerty", 
    Following {1}},
  { UserId: 3, 
    UserName:"Jaime Lannister", 
    Password:"qwerty", 
    Following {1,2,7}},
  {UserId: 4, 
    UserName:"Jon Snow", 
    Password:"qwerty", 
    Following {1,6,7}},
  {UserId: 5, 
    UserName:"Tyrion Lannister", 
    Password:"qwerty", 
    Following {1}},
  {UserId: 6, 
    UserName:"Daenerys Targaryen", 
    Password:"qwerty", 
    Following {4,7}},
  {UserId: 7, 
    UserName:"Cersei Lannister", 
    Password:"qwerty", 
    Following {3,6}}
  }

	http.HandleFunc("/", HomePage)
	// log.Fatal(http.ListenAndServe(":8080", nil))

	http.HandleFunc("/login", Login)
	err := http.ListenAndServe(":8080", nil) // setting listening port
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func HomePage(w http.ResponseWriter, r *http.Request) {

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

func Login(w http.ResponseWriter, r *http.Request) {
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

      //validate user
      for _, v := range userList {
        if v["UserName"] === userName && v["Password"] === password {
          //login success
          fmt.Fprintf(w, "Hello, %s. Welcom back.\n", userName)
        }

      }

      //login failed
      fmt.Fprintf(w, "Sorry, %s. Join us first.\n", userName)

	}
}
