package scrape

import (
	"fmt"
	"github.com/PuerkitoBio/goquery"
	//"io/ioutil"
	//"io"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"gopkg.in/mgo.v2"
)

type RoomInfo struct {
	Title    string `json:Title bson:Title`
	Cost     string `json:Cost bson:Cost`
	ImageUrl string `json:ImageUrl bson:ImageUrl`
}

var DBname = "spareroom"
var mongoUrl = "mongodb://egor2:qwer1234@ds153729.mlab.com:53729/spareroom"

var startUrl = "http://www.spareroom.co.uk/flatshare/search.pl?flatshare_type=offered&location_type=area&search="
var endUrl = "&miles_from_max=1&action=search&templateoveride=&show_results=&submit="

func ScrapeRoomsWithLocation(location string) ([]byte, error) {
	log.Println("Location for scrape: ", location)

	url := startUrl + location + endUrl
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatal(err)
	}
	var ErrorString = "Cant find such location, try another, or type it correct!"
	if doc.Find("#maincontent ul.listing-results article.panel-listing-result").Text() == "" {

		//return []byte(ErrorString), errors.New(ErrorString)
		return make([]byte, 0), errors.New(ErrorString)

	}

	mapRoomInfo := make([]RoomInfo, 11)

	session, err := mgo.Dial(mongoUrl)
	if err != nil {
		log.Println(err)
	}
	defer session.Close()

	doc.Find("#maincontent ul.listing-results article.panel-listing-result").Each(func(i int, s *goquery.Selection) {
		mapRoomInfo[i] = RoomInfo{
			Title:    s.Find("header.desktop a h1").Text(),
			Cost:     s.Find("strong.listingPrice").First().Text(),
			ImageUrl: s.Find("figure img").AttrOr("src", "No photo"),
		}
		RoomInfoColletion := session.DB(DBname).C("requestHistory")
		err = RoomInfoColletion.Insert(mapRoomInfo[i])
		if err != nil {
			log.Println(err)
		}

	})

	byteResults, err := json.Marshal(mapRoomInfo)
	if err != nil {
		log.Panicln(err)
	}

	return byteResults, nil
}

func TrialScrapeRooms(location string) ([]byte, error) {
	log.Println("Trial Location for scrape: ", location)
	log.Println(mongoUrl)

	url := startUrl + location + endUrl
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatal(err)
	}
	var ErrorString = "Cant find such location, try another, or type it correct!"
	if doc.Find("#maincontent ul.listing-results article.panel-listing-result").Text() == "" {

		//return []byte(ErrorString), errors.New(ErrorString)
		return make([]byte, 0), errors.New(ErrorString)

	}

	mapRoomInfo := make([]RoomInfo, 11)

	session, err := mgo.Dial(mongoUrl)
	if err != nil {
		log.Println(err)
	}
	defer session.Close()

	doc.Find("#maincontent ul.listing-results article.panel-listing-result").Each(func(i int, s *goquery.Selection) {
		mapRoomInfo[i] = RoomInfo{
			Title:    s.Find("header.desktop a h1").Text(),
			Cost:     s.Find("strong.listingPrice").First().Text(),
			ImageUrl: s.Find("figure img").AttrOr("src", "No photo"),
		}
		RoomInfoColletion := session.DB(DBname).C("requestHistory")
		err = RoomInfoColletion.Insert(mapRoomInfo[i])
		if err != nil {
			log.Println(err)
		}

	})

	byteResults, err := json.Marshal(mapRoomInfo)
	if err != nil {
		log.Panicln(err)
	}

	return byteResults, nil
}

func ScraperHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers",
		"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	r.ParseForm()

	location := r.FormValue("value")

	log.Println(r.FormValue("value"))
	data, err := ScrapeRoomsWithLocation(location)
	if err != nil {
		log.Println(err)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "Cant find such location, try another, or type it correct!")
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func TrialScraperHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers",
		"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

	location := "westminster"

	log.Println(location)
	data, err := TrialScrapeRooms(location)
	if err != nil {
		log.Println(err)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(w, "Cant find such location, try another, or type it correct!")
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}
