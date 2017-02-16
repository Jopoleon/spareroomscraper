package main

import (
	//"crypto/rand"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	//"crypto/sha1"
	//"io/ioutil"
	//"encoding/json"

	//"htmlhandlers"
	"scrape"
	"utils"

	"github.com/gorilla/context"
	"github.com/gorilla/sessions"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type UserInfo struct {
	Username      string `json:username bson:username`
	Password      string `json:password bson:password`
	Email         string `json:email bson:email`
	Loggedin      string `json:loggedin bson:loggedin`
	Registred     string `json:registered bson:registered`
	IsActivated   string `json:isActivated bson:isActivated`
	ActivationKey string `json:ActivationKey bson:ActivationKey`
}

var store = sessions.NewCookieStore([]byte("nRrHLlHcHH0u7fUz25Hje9m7uJ5SnJzP"))

var mongoUrl = "mongodb://egor_m:qwer1234@ds135029.mlab.com:35029/spareroom"

var DBname = "spareroom"

// func init() {
// 	dbsession, err := mgo.Dial(mongoUrl)
// 	if err != nil {
// 		log.Println(err)
// 	}
// }

func main() {
	ctl, err := NewController()
	if err != nil {
		log.Fatal(err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	fmt.Println("Server started on port: ", port)

	//Homepage
	http.HandleFunc("/", ctl.IndexHandler)
	//Scrape for rooms
	http.HandleFunc("/scrapelocation", scrape.ScraperHandler)
	//trial scrape for unregistered users
	http.HandleFunc("/trialscrapelocation", scrape.TrialScraperHandler)
	//sign up in system
	http.HandleFunc("/signup", SignUpHandler)
	//submit signup information
	http.HandleFunc("/signupsubmit", ctl.SignUpSubmitHandler)

	http.HandleFunc("/confirm", ctl.ConfirmSignUpHandler)

	http.HandleFunc("/login", LoginHandler)
	http.HandleFunc("/loginsubmit", ctl.LoginSubmitHandler)

	http.HandleFunc("/logout", ctl.LogoutSubmitHandler)

	//invoke static files(javascript, css, etc.)
	http.HandleFunc("/static/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, r.URL.Path[1:])
	})

	http.ListenAndServe(":"+port, context.ClearHandler(http.DefaultServeMux))

}

//https://www.andjosh.com/2015/01/31/middleware-in-go/
type Controller struct {
	// This will be our extensible type that will
	// be used as a common context type for our routes
	session *mgo.Session // our cloneable session
}

func NewController() (*Controller, error) {
	// This function will return to us a
	// Controller that has our common DB context.
	// We can then use it for multiple routes
	uri := mongoUrl
	if uri == "" {
		return nil, fmt.Errorf("no DB connection string provided")
	}
	session, err := mgo.Dial(uri)
	if err != nil {
		return nil, err
	}
	return &Controller{
		session: session,
	}, nil
}

func (ctl *Controller) SignUpSubmitHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "sessionRooms")
	if err != nil {
		log.Println(err)
		//http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	r.ParseForm()
	username := r.FormValue("username")
	password := r.FormValue("password")
	email := r.FormValue("email")
	currentUrl := r.FormValue("urlforactivation")
	log.Println(currentUrl)
	if !utils.ValidateEmail(email) {
		w.Write([]byte("Email incorrect "))
		return
	}
	if (username != "") && (password != "") && (email != "") {

		if !ctl.IsUserRegistered(username) {

			session.Values["registered"] = "true"
			session.Values["loggedin"] = "false"
			session.Values["username"] = username
			session.Values["password"] = password
			session.Values["email"] = email
			session.Values["isActivated"] = "false"
			newActivationKey := utils.GenerateKey32chars()
			session.Values["ActivationKey"] = newActivationKey
			session.Save(r, w)

			utils.SendEmailwithKey(newActivationKey, email, currentUrl)

			dbsession := ctl.session.Clone()
			defer dbsession.Close()

			RoomInfoColletion := dbsession.DB(DBname).C("usersInfo")
			err = RoomInfoColletion.Insert(
				&UserInfo{
					Registred:     session.Values["registered"].(string),
					Loggedin:      session.Values["loggedin"].(string),
					Username:      username, //from request
					Password:      password, //from request
					Email:         session.Values["email"].(string),
					IsActivated:   session.Values["isActivated"].(string),
					ActivationKey: newActivationKey,
				})
			if err != nil {
				log.Println(err)
			}
			w.Write([]byte("Registration successful! Check your email, and activate account"))
			//return
		} else {
			w.Write([]byte("User with such name or email is already exists "))
		}
	} else {
		w.Write([]byte("Some of registration fields are empty!"))
	}
	// log.Println(session.Values["password"], session.Values["username"])
}

func (ctl *Controller) LoginSubmitHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "sessionRooms")
	if err != nil {
		log.Println(err)
	}
	r.ParseForm()
	username := r.FormValue("username")
	password := r.FormValue("password")
	log.Print("Username login request: ", username, password)
	if ctl.IsUserRegistered(username) {

		if ctl.IsUserActivated(username) {

			dbsession := ctl.session.Clone()
			defer dbsession.Close()
			c := dbsession.DB(DBname).C("usersInfo")
			result := UserInfo{}
			err := c.Find(bson.M{"username": username}).One(&result)
			if err != nil {
				log.Println(err)
				w.Write([]byte("Username not found "))
				return
			}
			log.Printf("Data base userinfo: %+v \n", result)
			log.Print("pass from DB: ", result.Password, " pass from cookie: ", session.Values["password"])
			if result.Password == password {
				log.Println("inside cheking password")
				session.Values["loggedin"] = "true"
				session.Values["username"] = username
				session.Save(r, w)

				colQuerier := bson.M{"username": username}
				change := bson.M{"$set": bson.M{"loggedin": "true"}}
				err = c.Update(colQuerier, change)
				if err != nil {
					panic(err)
				}

				w.Write([]byte("You are logged!"))
				return
			} else {
				w.Write([]byte("Wrong password!"))
				return
			}
		} else {
			w.Write([]byte("Your account with username: " + username + " is not activated. Check your email: " + session.Values["email"].(string)))
			return

		}
	}
	// > db.usersInfo.distinct("username", {"registred":"true"})
	// > [ "egor", "egor2" ]

}

func (ctl *Controller) LogoutSubmitHandler(w http.ResponseWriter, r *http.Request) {
	session, err := store.Get(r, "sessionRooms")
	if err != nil {
		log.Println(err)
		//http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	log.Println("Inside Logout Handler, username for logout: ", session.Values["username"])
	if session.Values["username"].(string) != "" {
		username := session.Values["username"]
		if ctl.IsUserLogged(username.(string)) {

			dbsession := ctl.session.Clone()
			defer dbsession.Close()

			c := dbsession.DB(DBname).C("usersInfo")
			session.Values["loggedin"] = "false"
			session.Values["username"] = ""
			session.Save(r, w)
			colQuerier := bson.M{"username": username}
			change := bson.M{"$set": bson.M{"loggedin": "false"}}
			err = c.Update(colQuerier, change)
			if err != nil {
				panic(err)
			}
			w.Write([]byte("Succsesfuly logedout from " + username.(string)))
		} else {
			w.Write([]byte("You already have logged out"))
		}

	} else {
		w.Write([]byte("You dont have cookie session, please login first"))
	}

}

func (ctl *Controller) ConfirmSignUpHandler(w http.ResponseWriter, r *http.Request) {
	keyInUrl := r.URL.RawQuery
	dbsession := ctl.session.Clone()
	log.Println("Key from email link: ", keyInUrl)
	defer dbsession.Close()
	c := dbsession.DB(DBname).C("usersInfo")
	result := UserInfo{}

	err := c.Find(bson.M{"activationkey": keyInUrl}).One(&result)
	log.Println("Key from database: ", result.ActivationKey)
	if err != nil {
		log.Println(err)
		w.Write([]byte("Wrong activation key "))
		return
	}
	colQuerier := bson.M{"activationkey": keyInUrl}
	change := bson.M{"$set": bson.M{"isactivated": "true"}}
	err = c.Update(colQuerier, change)
	if err != nil {
		panic(err)
	}
	w.Write([]byte("Your account is active now"))
	return
}

func (ctl *Controller) IndexHandler(w http.ResponseWriter, r *http.Request) {

	session, err := store.Get(r, "sessionRooms")
	if err != nil {
		log.Println(err)
		//http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	log.Println("IndexHandler used")
	username := session.Values["username"]
	//log.Printf("Cookie values of current user: %+v \n", session.Values)
	log.Println("Current user: ", username)
	log.Println("Login status from DB: ", ctl.IsUserLogged(username.(string)))
	//session.Values["loggedin"] == "false" || session.Values["loggedin"] == nil ||
	if session.Values["loggedin"] == "false" || session.Values["loggedin"] == nil || !ctl.IsUserLogged(username.(string)) {
		session.Save(r, w)
		http.Redirect(w, r, "/login", 302)
		return
	} else {
		t, err := template.ParseFiles("static/index.html")
		if err != nil {
			fmt.Fprintf(w, err.Error())
		}
		t.ExecuteTemplate(w, "index.html", nil)
	}

}

func SignUpHandler(w http.ResponseWriter, r *http.Request) {

	t, err := template.ParseFiles("static/signup.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	t.ExecuteTemplate(w, "signup.html", nil)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {

	t, err := template.ParseFiles("static/login.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	t.ExecuteTemplate(w, "login.html", nil)

}

//check status ONLY in database!!!
func (ctl *Controller) IsUserRegistered(username string) bool {

	dbsession := ctl.session.Clone()
	defer dbsession.Close()

	c := dbsession.DB(DBname).C("usersInfo")
	result := UserInfo{}
	err := c.Find(bson.M{"username": username}).One(&result)
	if err != nil {
		log.Println(err)

		return false
	}
	if result.Registred == "true" {
		return true
	} else {

		return false
	}
}
func (ctl *Controller) IsUserLogged(username string) bool {

	dbsession := ctl.session.Clone()
	defer dbsession.Close()

	c := dbsession.DB(DBname).C("usersInfo")
	result := UserInfo{}
	err := c.Find(bson.M{"username": username}).One(&result)
	if err != nil {
		log.Println(err)
		return false
	}
	if result.Loggedin == "true" {
		return true
	} else {
		return false
	}
}
func (ctl *Controller) IsUserActivated(username string) bool {

	dbsession := ctl.session.Clone()
	defer dbsession.Close()

	c := dbsession.DB(DBname).C("usersInfo")
	result := UserInfo{}
	err := c.Find(bson.M{"username": username}).One(&result)
	if err != nil {
		log.Println(err)
		return false
	}
	if result.IsActivated == "true" {
		return true
	} else {
		return false
	}
}
