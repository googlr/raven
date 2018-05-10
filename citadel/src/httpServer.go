package main

import (
	// "errors"
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"net/rpc"
	"strings"
	"unicode"
)

var (
	// back end server address
	serverAddress string = "localhost:"
	port          string = "9527"
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
	IsFollowing bool
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

type SendMessageReply struct {
	SendMsgStatus  bool
	SendMsgProfile UserProfile
}

func init() {
	gob.Register(&UserProfile{})
}

var userFollowingMap = map[string]map[string]string{
	"Ned.Stark@Winterfell.com": {
		"Robert.Baratheon@kingslanding.com": "Robert.Baratheon@kingslanding.com",
	},
	"Robert.Baratheon@kingslanding.com": {
		"Ned.Stark@Winterfell.com": "Ned.Stark@Winterfell.com",
	},
	"Jaime.Lannister@CasterlyRock.com": {
		"Tyrion.Lannister@CasterlyRock.com": "Tyrion.Lannister@CasterlyRock.com",
		"Cersei.Lannister@CasterlyRock.com": "Cersei.Lannister@CasterlyRock.com",
	},
	"Jon.Snow@Winterfell.com": {
		"Ned.Stark@Winterfell.com":           "Ned.Stark@Winterfell.com",
		"Daenerys.Targaryen@Dragonstone.com": "Daenerys.Targaryen@Dragonstone.com",
	},
	"Tyrion.Lannister@CasterlyRock.com": {
		"Ned.Stark@Winterfell.com":           "Ned.Stark@Winterfell.com",
		"Robert.Baratheon@kingslanding.com":  "Robert.Baratheon@kingslanding.com",
		"Jaime.Lannister@CasterlyRock.com":   "Jaime.Lannister@CasterlyRock.com",
		"Jon.Snow@Winterfell.com":            "Jon.Snow@Winterfell.com",
		"Daenerys.Targaryen@Dragonstone.com": "Daenerys.Targaryen@Dragonstone.com",
		"Cersei.Lannister@CasterlyRock.com":  "Cersei.Lannister@CasterlyRock.com",
	},
	"Daenerys.Targaryen@Dragonstone.com": {
		"Tyrion.Lannister@CasterlyRock.com": "Tyrion.Lannister@CasterlyRock.com",
		"Jon.Snow@Winterfell.com":           "Jon.Snow@Winterfell.com",
		"Cersei.Lannister@CasterlyRock.com": "Cersei.Lannister@CasterlyRock.com",
	},
	"Cersei.Lannister@CasterlyRock.com": {
		"Jaime.Lannister@CasterlyRock.com": "Jaime.Lannister@CasterlyRock.com",
	},
}

var userFollowerMap = map[string]map[string]string{
	"Ned.Stark@Winterfell.com": {
		"Robert.Baratheon@kingslanding.com": "Robert.Baratheon@kingslanding.com",
		"Jon.Snow@Winterfell.com":           "Jon.Snow@Winterfell.com",
		"Tyrion.Lannister@CasterlyRock.com": "Tyrion.Lannister@CasterlyRock.com",
	},

	"Robert.Baratheon@kingslanding.com": {
		"Ned.Stark@Winterfell.com":          "Ned.Stark@Winterfell.com",
		"Tyrion.Lannister@CasterlyRock.com": "Tyrion.Lannister@CasterlyRock.com",
	},
	"Jaime.Lannister@CasterlyRock.com": {
		"Tyrion.Lannister@CasterlyRock.com": "Tyrion.Lannister@CasterlyRock.com",
		"Cersei.Lannister@CasterlyRock.com": "Cersei.Lannister@CasterlyRock.com",
	},
	"Jon.Snow@Winterfell.com": {
		"Tyrion.Lannister@CasterlyRock.com":  "Tyrion.Lannister@CasterlyRock.com",
		"Daenerys.Targaryen@Dragonstone.com": "Daenerys.Targaryen@Dragonstone.com",
	},
	"Tyrion.Lannister@CasterlyRock.com": {
		"Jaime.Lannister@CasterlyRock.com":   "Jaime.Lannister@CasterlyRock.com",
		"Daenerys.Targaryen@Dragonstone.com": "Daenerys.Targaryen@Dragonstone.com",
	},
	"Daenerys.Targaryen@Dragonstone.com": {
		"Jon.Snow@Winterfell.com":           "Jon.Snow@Winterfell.com",
		"Tyrion.Lannister@CasterlyRock.com": "Tyrion.Lannister@CasterlyRock.com",
	},
	"Cersei.Lannister@CasterlyRock.com": {
		"Jaime.Lannister@CasterlyRock.com":   "Jaime.Lannister@CasterlyRock.com",
		"Tyrion.Lannister@CasterlyRock.com":  "Tyrion.Lannister@CasterlyRock.com",
		"Daenerys.Targaryen@Dragonstone.com": "Daenerys.Targaryen@Dragonstone.com",
	},
}

var userMsgMap = map[string][]Message{
	"Ned.Stark@Winterfell.com": []Message{
		{
			SenderEmail: "Ned.Stark@Winterfell.com",
			SenderName:  "Ned Stark",
			MsgContent:  "The winters are hard but the Starks will endure. We always have."},
		{
			SenderEmail: "Ned.Stark@Winterfell.com",
			SenderName:  "Ned Stark",
			MsgContent:  "The next time we see each other, we'll talk about your mother. I promise."},
	},
	"Robert.Baratheon@kingslanding.com": []Message{
		{
			SenderEmail: "Robert.Baratheon@kingslanding.com",
			SenderName:  "Robert Baratheon",
			MsgContent:  "I'm not trying to honor you. I'm trying to get you to run my kingdom while I eat, drink, and whore my way to an early grave."},
		{
			SenderEmail: "Robert.Baratheon@kingslanding.com",
			SenderName:  "Robert Baratheon",
			MsgContent:  "Kill the F**cing Boar."},
	},
	"Jaime.Lannister@CasterlyRock.com": []Message{
		{
			SenderEmail: "Jaime.Lannister@CasterlyRock.com",
			SenderName:  "Jaime Lannister",
			MsgContent:  "The things I do for love."},
		{
			SenderEmail: "Jaime.Lannister@CasterlyRock.com",
			SenderName:  "Jaime Lannister",
			MsgContent:  "There are no men like me. Only me."},
	},
	"Jon.Snow@Winterfell.com": []Message{
		{
			SenderEmail: "Jon.Snow@Winterfell.com",
			SenderName:  "Jon Snow",
			MsgContent:  "I am not a Stark."},
		{
			SenderEmail: "Jon.Snow@Winterfell.com",
			SenderName:  "Jon Snow",
			MsgContent:  "My watch is ended."},
	},
	"Tyrion.Lannister@CasterlyRock.com": []Message{
		{
			SenderEmail: "Tyrion.Lannister@CasterlyRock.com",
			SenderName:  "Tyrion Lannister",
			MsgContent:  "Never forget what you are, the rest of the world will not. Wear it like armor and it can never be used to hurt you."},
		{
			SenderEmail: "Tyrion.Lannister@CasterlyRock.com",
			SenderName:  "Tyrion Lannister",
			MsgContent:  "I have to disagree. Death is so final, yet life is full of possibilities."},
		{
			SenderEmail: "Tyrion.Lannister@CasterlyRock.com",
			SenderName:  "Tyrion Lannister",
			MsgContent:  "A mind needs books like a sword needs a whetstone."},
		{
			SenderEmail: "Tyrion.Lannister@CasterlyRock.com",
			SenderName:  "Tyrion Lannister",
			MsgContent:  "I have a tender spot in my heart for cripples and bastards and broken things."},
		{
			SenderEmail: "Tyrion.Lannister@CasterlyRock.com",
			SenderName:  "Tyrion Lannister",
			MsgContent:  "We’ve had vicious kings and we’ve had idiot kings, but I don’t know if we’ve ever been cursed with a vicious idiot boy king!"},
		{
			SenderEmail: "Tyrion.Lannister@CasterlyRock.com",
			SenderName:  "Tyrion Lannister",
			MsgContent:  "It’s hard to put a leash on a dog once you’ve put a crown on its head."},
		{
			SenderEmail: "Tyrion.Lannister@CasterlyRock.com",
			SenderName:  "Tyrion Lannister",
			MsgContent:  "Drinking and lust. No man can match me in these things. I am the god of tits and wine… I shall build a shrine to myself at the next brothel I visit."},
		{
			SenderEmail: "Tyrion.Lannister@CasterlyRock.com",
			SenderName:  "Tyrion Lannister",
			MsgContent:  "Every time we deal with an enemy, we create two more."},
	},
	"Daenerys.Targaryen@Dragonstone.com": []Message{
		{
			SenderEmail: "Daenerys.Targaryen@Dragonstone.com",
			SenderName:  "Daenerys Targaryen",
			MsgContent:  "I am the blood of the dragon. I must be strong. I must have fire in my eyes when I face them, not tears."},
		{
			SenderEmail: "Daenerys.Targaryen@Dragonstone.com",
			SenderName:  "Daenerys Targaryen",
			MsgContent:  "Valar morghulis. Daenerys Targaryen: Yes. All men must die, but we are not men."},
	},
	"Cersei.Lannister@CasterlyRock.com": []Message{
		{
			SenderEmail: "Cersei.Lannister@CasterlyRock.com",
			SenderName:  "Cersei Lannister",
			MsgContent:  "Everyone who isn’t us is an enemy."},
		{
			SenderEmail: "Cersei.Lannister@CasterlyRock.com",
			SenderName:  "Cersei Lannister",
			MsgContent:  "When you play the game of thrones you win or you die. There is no middle ground."},
		{
			SenderEmail: "Cersei.Lannister@CasterlyRock.com",
			SenderName:  "Cersei Lannister",
			MsgContent:  "I choose violence."},
	},
}

var userProfileMap = map[string]UserProfile{
	//"Ned Stark"
	"Ned.Stark@Winterfell.com": UserProfile{
		UserEmail:    "Ned.Stark@Winterfell.com",
		UserId:       1,
		UserName:     "Ned Stark",
		UserBio:      "Winter is coming.",
		FollowingNum: 1,
		FollowerNum:  3,
	},
	//"Robert Baratheon"
	"Robert.Baratheon@Kingslanding.com": UserProfile{
		UserEmail:    "Robert.Baratheon@Kingslanding.com",
		UserId:       2,
		UserName:     "Robert Baratheon",
		UserBio:      "Ours Is the Fury.",
		FollowingNum: 1,
		FollowerNum:  2,
	},
	//"Jaime Lannister"
	"Jaime.Lannister@CasterlyRock.com": UserProfile{
		UserEmail:    "Jaime.Lannister@CasterlyRock.com",
		UserId:       3,
		UserName:     "Jaime Lannister",
		UserBio:      "A Lannister Always Pays His Debts.",
		FollowingNum: 2,
		FollowerNum:  2,
	},
	//"Jon Snow"
	"Jon.Snow@Winterfell.com": UserProfile{
		UserEmail:    "Jon.Snow@Winterfell.com",
		UserId:       4,
		UserName:     "Jon Snow",
		UserBio:      "Winter is coming. Meanwhile, I am not a Stark.",
		FollowingNum: 2,
		FollowerNum:  2,
	},
	//"Tyrion Lannister"
	"Tyrion.Lannister@CasterlyRock.com": UserProfile{
		UserEmail:    "Tyrion.Lannister@CasterlyRock.com",
		UserId:       5,
		UserName:     "Tyrion Lannister",
		UserBio:      "A Lannister Always Pays His Debts.",
		FollowingNum: 6,
		FollowerNum:  2,
	},
	"Daenerys.Targaryen@Dragonstone.com": UserProfile{
		UserEmail:    "Daenerys.Targaryen@Dragonstone.com",
		UserId:       6,
		UserName:     "Daenerys Targaryen",
		UserBio:      "Fire and Blood.",
		FollowingNum: 3,
		FollowerNum:  2,
	},
	"Cersei.Lannister@CasterlyRock.com": UserProfile{
		UserEmail:    "Cersei.Lannister@CasterlyRock.com",
		UserId:       7,
		UserName:     "Cersei Lannister",
		UserBio:      "A Lannister Always Pays His Debts.",
		FollowingNum: 1,
		FollowerNum:  3,
	},
}

var userCredentialMap = map[string]UserCredential{
	"Ned.Stark@Winterfell.com": UserCredential{
		UserEmail: "Ned.Stark@Winterfell.com",
		Password:  "qwerty",
	},
	"Robert.Baratheon@Kingslanding.com": UserCredential{
		UserEmail: "Robert.Baratheon@Kingslanding.com",
		Password:  "qwerty",
	},
	"Jaime.Lannister@CasterlyRock.com": UserCredential{
		UserEmail: "Jaime.Lannister@CasterlyRock.com",
		Password:  "qwerty",
	},
	"Jon.Snow@Winterfell.com": UserCredential{
		UserEmail: "Jon.Snow@Winterfell.com",
		Password:  "qwerty",
	},
	"Tyrion.Lannister@CasterlyRock.com": UserCredential{
		UserEmail: "Tyrion.Lannister@CasterlyRock.com",
		Password:  "qwerty",
	},
	"Daenerys.Targaryen@Dragonstone.com": UserCredential{
		UserEmail: "Daenerys.Targaryen@Dragonstone.com",
		Password:  "qwerty",
	},
	"Cersei.Lannister@CasterlyRock.com": UserCredential{
		UserEmail: "Cersei.Lannister@CasterlyRock.com",
		Password:  "qwerty",
	},
}

func (usr *UserProfile) UserLoginHandler(args *LoginArgs, reply *LoginReply) error {
	//validation
	userEmail := args.UserLoginEmail
	userPswd := args.UserLoginPassword
	userProfile, ok := userProfileMap[userEmail]
	if ok == true {
		userCredentialStored, ok2 := userCredentialMap[userEmail]
		if ok2 == true {
			if userCredentialStored.Password == userPswd {
				fmt.Println("UserLogin: success.")
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
	}
	reply.UserLoginStatus = false
	reply.UserLoginProfile = UserProfile{
		UserEmail:    "",
		UserId:       0,
		UserName:     "",
		UserBio:      "",
		FollowingNum: 0,
		FollowerNum:  0,
	}
	return nil
}

func (usr *UserProfile) UserSignUpHandler(args *SignUpArgs, reply *SignUpReply) error {
	// UserSignUpEmail, UserSignUpName, UserSignUpPassword
	userEmail := args.UserSignUpEmail
	userName := args.UserSignUpName
	//validation

	_, ok := userProfileMap[userEmail]
	if ok == true {
		fmt.Println("UserSignUp: userName already exist, be more creative.")
		reply.UserSignUpStatus = false

		reply.UserSignUpProfile = UserProfile{
			UserEmail:    "",
			UserId:       0,
			UserName:     "",
			UserBio:      "",
			FollowingNum: 0,
			FollowerNum:  0,
		}
		return nil
	} else {
		fmt.Println("UserSignUp: success.")
		reply.UserSignUpStatus = true

		//create a nil User with Profile and Credential
		signUpUserProfile := UserProfile{
			UserEmail:    userEmail,
			UserId:       len(userProfileMap) + 1,
			UserName:     userName,
			UserBio:      "",
			FollowingNum: 0,
			FollowerNum:  0,
		}

		signUpUserCredential := UserCredential{
			UserEmail: userEmail,
			Password:  args.UserSignUpPassword,
		}
		//add new User to database
		userProfileMap[userEmail] = signUpUserProfile
		userCredentialMap[userEmail] = signUpUserCredential

		reply.UserSignUpProfile = signUpUserProfile

		//add Log
		// log := "Inset new user:"
		// addLog(&log)
		return nil
	}

	// Never come here
	return nil
}

func (usr *UserProfile) SendMessageHandler(args *Message, reply *SendMessageReply) error {
	// SenderEmail string, SenderName  string, MsgContent string
	//validation
	userEmail := args.SenderEmail
	_, ok := userProfileMap[userEmail]
	if ok == true {
		fmt.Println("SendMessage: success.")
		reply.SendMsgStatus = true
		msgMap, ok := userMsgMap[userEmail]
		if ok == false {
			//This user has never post msg before
			userMsgMap[userEmail] = []Message{*args}
		}
		msgMap = append(msgMap, *args)

		userMsgMap[userEmail] = msgMap

		reply.SendMsgProfile = userProfileMap[userEmail] //!! re extract user profile
		return nil
	} else {
		fmt.Println("SendMessage: userEmail does not exist, go to sign up.")
		reply.SendMsgStatus = false
		reply.SendMsgProfile = UserProfile{
			UserEmail:    "",
			UserId:       0,
			UserName:     "",
			UserBio:      "",
			FollowingNum: 0,
			FollowerNum:  0,
		}
		return nil
	}
	// Never come here
	return nil
}

func (usr *UserProfile) GetUserMsgFromFollowingHandler(args *UserProfile, reply *[]Message) error {

	// SenderEmail string, SenderName  string, MsgContent string
	//validation
	userEmail := args.UserEmail
	currFollowingMap, ok := userFollowingMap[userEmail]
	if ok == true {
		// get Following UserEmail
		fmt.Println("GetUserMsgFromFollowing: success.")
		for _, v := range currFollowingMap {
			// get Msg for Each UserEmail
			msgArray, ok := userMsgMap[v]
			if ok == true {
				for _, msg := range msgArray {
					*reply = append(*reply, msg)
				}
			} else {
				//
				fmt.Println("GetUserMsgFromFollowing: user has no msg post.")
			}
		}
		return nil
	} else {
		fmt.Println("GetUserMsgFromFollowing: userEmail does not exist, go to sign up.")
		reply = nil
		return nil
	}
	return nil
}

func (usr *UserProfile) GetSearchResultsHandler(args *string, reply *[]UserProfile) error {
	// Make a Regex to say we only want
	// reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// searchKeyword := reg.ReplaceAllString((*args), "")

	searchKeyword := strings.ToLower(
		strings.TrimFunc((*args), func(r rune) bool {
			return !unicode.IsLetter(r) && !unicode.IsNumber(r)
		}))

	// search the userProfileMap
	fmt.Printf("GetSearchResults: searching %s.\n", searchKeyword)
	for _, v := range userProfileMap {
		userNameLower := strings.ToLower(
			strings.TrimFunc(v.UserName, func(r rune) bool {
				return !unicode.IsLetter(r) && !unicode.IsNumber(r)
			}))
		// fmt.Printf("GetSearchResults: matching %s.\n", userNameLower)
		if strings.Contains(userNameLower, searchKeyword) {
			*reply = append(*reply, v)
		}

	}

	fmt.Println("GetSearchResults: success.")
	return nil
}

func (usr *UserProfile) GetUserProfileHandler(args *GetProfilePageArgs, reply *GetProfilePageReply) error {
	fmt.Println("Getting User Profile.")
	userEmail := (*args).User
	currUserEmail := (*args).CurrentUser
	fmt.Println(userEmail)
	fmt.Println(currUserEmail)

	currFollowingMap, ok := userFollowingMap[currUserEmail]
	if ok == true {
		_, ok2 := currFollowingMap[userEmail]
		if ok2 == true {
			//current user has followed the user
			fmt.Println("Current user is following the user.")
			reply.IsFollowing = true
			reply.User = userProfileMap[userEmail]
			fmt.Println(reply)
			return nil
		} else {
			//current user has not followed the user
			fmt.Println("Current user has not followed the user.")
			reply.IsFollowing = false
			reply.User = userProfileMap[userEmail]
		}

	} else {
		// current user has never followed anyone, is empty in following map
		fmt.Println("Current user has not followed anyone.")
		reply.IsFollowing = false
		reply.User = userProfileMap[userEmail]
	}
	return nil
}

func addLog(args *string) {
	client, err := rpc.DialHTTP("tcp", serverAddress+port)
	if err != nil {
		log.Fatal("dialing:", err)
	}
	// Synchronous call
	// RPC call
	var logReply bool
	err = client.Call("UserProfile.AddLogHandler", args, &logReply)
	if err != nil {
		log.Fatal("User Login error:", err)
	}
	if logReply == true {
		fmt.Println("addLog: success.")
	} else {
		fmt.Println("addLog: failed.")
	}
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
