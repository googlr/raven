package main

import (
	// "errors"
	"fmt"
	"net/http"
	"net/rpc"
)

// type UserAccount struct {
// 	UserName string
// 	Password string
// }

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

type Arith int

type LoginArgs struct {
	UserLoginName, UserLoginPassword string
}

type LoginReply struct {
	UserLoginStatus  bool
	UserLoginProfile User
}

func (usr *User) UserLoginValidation(args *LoginArgs, reply *LoginReply) error {
	//validation
	for _, v := range userList {
		// fmt.Println("User: ", v.UserName)

		if v.UserName == args.UserLoginName {
			// && v.Password == password[0]
			// 	//login success
			reply.UserLoginStatus = true
			reply.UserLoginProfile = v
			return nil
		}
	}
	reply.UserLoginStatus = false
	reply.UserLoginProfile = User{UserId: 0,
		UserName:  "",
		Password:  "",
		Following: []int{}}
	return nil
}

func (usr *User) UserSignUpHandler(args *LoginArgs, reply *LoginReply) error {
	//validation
	for _, v := range userList {
		// fmt.Println("User: ", v.UserName)

		if v.UserName == args.UserLoginName {
			// && v.Password == password[0]
			// userName already registered!
			reply.UserLoginStatus = false
			reply.UserLoginProfile = User{UserId: 0,
				UserName:  "",
				Password:  "",
				Following: []int{},
			}
			return nil
		}
	}

	reply.UserLoginStatus = true
	reply.UserLoginProfile = User{
		UserId:    len(userList) + 1,
		UserName:  args.UserLoginName,
		Password:  "qwerty",
		Following: []int{},
	}
	return nil
}

func main() {

	user := new(User)
	rpc.Register(user)
	rpc.HandleHTTP()

	err := http.ListenAndServe(":1234", nil)
	if err != nil {
		fmt.Println(err.Error())
	}

}
