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

type UserProfile struct {
	UserId   int
	UserName string
	// Password  string
	Following []int
	PostMsg   []Message
}

type UserCredential struct {
	UserId   int
	Password string
}

type Message struct {
	SenderId int
	// timeStamp Time
	Content string
}

var userProfileMap = map[string]UserProfile{
	"Ned Stark": UserProfile{UserId: 1,
		UserName: "Ned Stark",
		// Password:  "qwerty",
		Following: []int{2},
		PostMsg: []Message{
			{
				SenderId: 1,
				Content:  "The winters are hard but the Starks will endure. We always have."},
			{
				SenderId: 1,
				Content:  "The next time we see each other, we'll talk about your mother. I promise."}}},
	"Robert Baratheon": UserProfile{UserId: 2,
		UserName: "Robert Baratheon",
		// Password:  "qwerty",
		Following: []int{1},
		PostMsg: []Message{
			{
				SenderId: 2,
				Content:  "I'm not trying to honor you. I'm trying to get you to run my kingdom while I eat, drink, and whore my way to an early grave."}}},
	"Jaime Lannister": UserProfile{UserId: 3,
		UserName: "Jaime Lannister",
		// Password:  "qwerty",
		Following: []int{1, 2, 7},
		PostMsg: []Message{
			{
				SenderId: 3,
				Content:  "The things I do for love."}}},
	"Jon Snow": UserProfile{UserId: 4,
		UserName: "Jon Snow",
		// Password:  "qwerty",
		Following: []int{1, 6, 7},
		PostMsg: []Message{
			{
				SenderId: 4,
				Content:  "I am not a Stark."},
			{
				SenderId: 4,
				Content:  "My watch is ended."}}},
	"Tyrion Lannister": UserProfile{UserId: 5,
		UserName: "Tyrion Lannister",
		// Password:  "qwerty",
		Following: []int{1},
		PostMsg: []Message{
			{
				SenderId: 5,
				Content:  "Never forget what you are, the rest of the world will not. Wear it like armor and it can never be used to hurt you."},
			{
				SenderId: 5,
				Content:  "I have to disagree. Death is so final, yet life is full of possibilities."},
			{
				SenderId: 5,
				Content:  "A mind needs books like a sword needs a whetstone."},
			{
				SenderId: 5,
				Content:  "I have a tender spot in my heart for cripples and bastards and broken things."},
			{
				SenderId: 5,
				Content:  "We’ve had vicious kings and we’ve had idiot kings, but I don’t know if we’ve ever been cursed with a vicious idiot boy king!"},
			{
				SenderId: 5,
				Content:  "It’s hard to put a leash on a dog once you’ve put a crown on its head."},
			{
				SenderId: 5,
				Content:  "Drinking and lust. No man can match me in these things. I am the god of tits and wine… I shall build a shrine to myself at the next brothel I visit."},
			{
				SenderId: 5,
				Content:  "Every time we deal with an enemy, we create two more."}}},
	"Daenerys Targaryen": UserProfile{UserId: 6,
		UserName: "Daenerys Targaryen",
		// Password:  "qwerty",
		Following: []int{4, 7},
		PostMsg: []Message{
			{
				SenderId: 6,
				Content:  "I am the blood of the dragon. I must be strong. I must have fire in my eyes when I face them, not tears."},
			{
				SenderId: 6,
				Content:  "Valar morghulis. Daenerys Targaryen: Yes. All men must die, but we are not men."}}},
	"Cersei Lannister": UserProfile{UserId: 7,
		UserName: "Cersei Lannister",
		// Password:  "qwerty",
		Following: []int{3, 6},
		PostMsg: []Message{
			{
				SenderId: 7,
				Content:  "Everyone who isn’t us is an enemy."},
			{
				SenderId: 7,
				Content:  "When you play the game of thrones you win or you die. There is no middle ground."},
			{
				SenderId: 7,
				Content:  "I choose violence."}}}}

var userCredentialMap = map[string]UserCredential{
	"Ned Stark": UserCredential{UserId: 1,
		Password: "qwerty",
	},
	"Robert Baratheon": UserCredential{UserId: 1,
		Password: "qwerty",
	},
	"Jaime Lannister": UserCredential{UserId: 1,
		Password: "qwerty",
	},
	"Jon Snow": UserCredential{UserId: 1,
		Password: "qwerty",
	},
	"Tyrion Lannister": UserCredential{UserId: 1,
		Password: "qwerty",
	},
	"Daenerys Targaryen": UserCredential{UserId: 1,
		Password: "qwerty",
	},
	"Cersei Lannister": UserCredential{UserId: 1,
		Password: "qwerty",
	},
}

type LoginArgs struct {
	UserLoginName, UserLoginPassword string
}

type LoginReply struct {
	UserLoginStatus  bool
	UserLoginProfile UserProfile
}

func (usr *UserProfile) UserLoginValidation(args *LoginArgs, reply *LoginReply) error {
	//validation
	userName := args.UserLoginName
	userPswd := args.UserLoginPassword
	userProfile, ok := userProfileMap[userName]
	if ok == true {
		userCredentialStored, ok2 := userCredentialMap[userName]
		if ok2 == true {
			if userCredentialStored.Password == userPswd {
				reply.UserLoginStatus = true
				reply.UserLoginProfile = userProfile
				return nil
			} else {
				fmt.Println("UserLogin: password not match.")
			}
		} else {
			fmt.Println("UserLogin: user password not found on server, this should not happen.")
		}
	} else {
		fmt.Println("UserLogin: userName does not exist, go to sign up.")
		reply.UserLoginStatus = false
		reply.UserLoginProfile = UserProfile{
			UserId:   0,
			UserName: "",
			// Password:  "",
			Following: []int{},
			PostMsg:   []Message{}}
		return nil
	}
	// Never come here
	return nil
}

func (usr *UserProfile) UserSignUpHandler(args *LoginArgs, reply *LoginReply) error {
	userName := args.UserLoginName
	//validation

	_, ok := userProfileMap[userName]
	if ok == true {
		fmt.Println("UserSignUp: userName already exist, be more creative.")
		reply.UserLoginStatus = false

		reply.UserLoginProfile = UserProfile{
			UserId:   0,
			UserName: "",
			// Password:  "",
			Following: []int{},
			PostMsg:   []Message{}}
		return nil
	} else {
		fmt.Println("UserSignUp: success.")
		reply.UserLoginStatus = true

		//create a new User with Profile and Credential
		signUpUserProfile := UserProfile{
			UserId:   len(userProfileMap) + 1,
			UserName: args.UserLoginName,
			// Password:  args.UserLoginPassword,
			Following: []int{},
			PostMsg:   []Message{}}

		signUpUserCredential := UserCredential{
			UserId:   signUpUserProfile.UserId,
			Password: args.UserLoginPassword}
		//add new User to database
		userProfileMap[userName] = signUpUserProfile
		userCredentialMap[userName] = signUpUserCredential

		reply.UserLoginProfile = signUpUserProfile
		return nil
	}

	// Never come here
	return nil
}

func main() {

	user := new(UserProfile)
	rpc.Register(user)
	rpc.HandleHTTP()

	err := http.ListenAndServe(":1234", nil)
	if err != nil {
		fmt.Println(err.Error())
	}

}
