package utils

import (
	"crypto/rand"
	"fmt"
	"net/smtp"
	//"github.com/PuerkitoBio/goquery"
	//"io/ioutil"
	//"io"
	//"encoding/json"
	//"errors"
	"log"
	"strconv"
	"time"

	"regexp"
)

var url_start = "http://www.spareroom.co.uk"
var searchEndpoint = "/flatshare/search.pl?"
var full_request_url = "http://www.spareroom.co.uk/flatshare/search.pl?flatshare_type={offered|wanted|buddyup}&location_type=area&search={searchterm}&miles_from_max={miles}&showme_rooms=Y&showme_1beds=Y&showme_buddyup_properties=Y&min_rent={mincost}&max_rent={maxcost}&per=pcm&no_of_rooms=&min_term=0&max_term=0&available_search=N&day_avail=&mon_avail=&year_avail=&min_age_req=&max_age_req=&min_beds=&max_beds=&keyword=&searchtype=advanced%20&editing=&mode=&nmsq_mode=&action=search&templateoveride=&show_results=&submit="
var url_end = "&action=search&templateoveride=&show_results=&submit="

var firtsTryUrlsReq = "http://www.spareroom.co.uk/flatshare/search.pl?flatshare_type=offered&location_type=area&search=westminster&miles_from_max=1&action=search&templateoveride=&show_results=&submit="
var loca = "belgravia"

var startUrl = "http://www.spareroom.co.uk/flatshare/search.pl?flatshare_type=offered&location_type=area&search="
var endUrl = "&miles_from_max=1&action=search&templateoveride=&show_results=&submit="

func GenerateKey32chars() string {
	buf := make([]byte, 16)
	_, err := rand.Read(buf)
	if err != nil {
		panic(err) // out of randomness, should never happen
	}
	return fmt.Sprintf("%x", buf)
	// or hex.EncodeToString(buf)
	// or base64.StdEncoding.EncodeToString(buf)
}

//https://www.socketloop.com/tutorials/golang-validate-email-address-with-regular-expression
func ValidateEmail(email string) bool {
	Re := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	return Re.MatchString(email)
}

func SendEmailwithKey(key, addres, currentActivationUrl string) {
	fmt.Println("Start time: ", time.Now())
	//Sendkey2017
	type EmailUser struct {
		Username    string
		Password    string
		EmailServer string
		Port        int
	}

	emailUser := &EmailUser{"spareroommailserver", "Sendkey2017", "smtp.gmail.com", 587}

	auth := smtp.PlainAuth("",
		emailUser.Username,
		emailUser.Password,
		emailUser.EmailServer,
	)

	var err error

	link := currentActivationUrl + "?" + key
	//linkTag := "<a href='" + link + "'>" + link + "</a>"
	msg := []byte("To: " + addres + "\r\n" +
		"Subject: Activation letter from SpareRoomScraper\r\n" +
		"\r\n" +
		"This is your activation link: \r\n" + link)
	//<a href="http://localhost:3000/confirm?1d27bcb9ab2297bc8c569ffd8adc2902">https://localhost:3000/confirm?1d27bcb9ab2297bc8c569ffd8adc2902
	err = smtp.SendMail(emailUser.EmailServer+":"+strconv.Itoa(emailUser.Port),
		auth,
		emailUser.Username,
		[]string{addres},
		msg)
	if err != nil {
		log.Print("ERROR: attempting to send a mail ", err)
	}
	fmt.Println("End time: ", time.Now())
}

// Old versions of scraper are presented below
//
// func ScrapeRoomsWithLocation(location string) ([]byte, error) {
// 	log.Println("Location for scrape: ", location)

// 	url := startUrl + location + endUrl
// 	doc, err := goquery.NewDocument(url)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	if doc.Find("#maincontent ul.listing-results article.panel-listing-result").Text() == "" {
// 		return make([]byte, 0), errors.New("Cant find such location, try another, or type it correct!")
// 		// callback = append(callback, "Cant find such location, try another, or type it correct!")
// 		// return callback
// 	}
// 	//var searchResults []RoomInfo
// 	var byteResults []byte
// 	//roomsMap = make(map[string]string)
// 	//#maincontent > ul > li:nth-child(1) > article > figure > a > img

// 	doc.Find("#maincontent ul.listing-results article.panel-listing-result").Each(func(i int, s *goquery.Selection) {
// 		roomsMap := map[string]string{
// 			"Title":    s.Find("header.desktop a h1").Text(),
// 			"Cost":     s.Find("strong.listingPrice").First().Text(),
// 			"ImageUrl": s.Find("figure img").AttrOr("src", "No photo"),
// 		}
// 		log.Println("Data BEFORE marshaling json: ", roomsMap)
// 		//fmt.Printf("%s", roomsMap)
// 		//fmt.Println(len(`{"Cost":"£800 pcm","ImageUrl":"//photos2.spareroom.co.uk/images/flatshare/listings/cw100h100/27/25/27250129.jpg","Title":"*Double room in central location, Edgware Road"}`))

// 		// ri := RoomInfo{}

// 		// ri.Title = s.Find("header.desktop a h1").Text()
// 		// ri.Cost = s.Find("strong.listingPrice").First().Text()
// 		// ri.ImageUrl = s.Find("figure img").AttrOr("src", "No photo")

// 		// a := []byte("Hello")
// 		// b := []byte(" World!")
// 		// c := [][]byte{a, b}

// 		// d := []byte{}

// 		// d = bytes.Join(c, []byte(", "))
// 		// fmt.Println(string(d))

// 		// e := []byte{}

// 		// e = bytes.Join(c, []byte("--"))
// 		// fmt.Println(string(e))

// 		b, err := json.Marshal(roomsMap)
// 		if err != nil {
// 			log.Panicln(err)
// 		}
// 		byteResults = append(byteResults, append(b, []byte(", ")...)...)
// 		// if err != nil {
// 		// 	return make([]RoomInfo, 0), err.Error())
// 		// }
// 		// title := s.Find("header.desktop a h1").Text()
// 		// weekCost := s.Find("strong.listingPrice").First().Text()
// 		// ImageUrl := s.Find("figure img").AttrOr("src", "No photo")
// 		//searchResults = append(searchResults, ri)

// 		//fmt.Println(title + " " + weekCost + " " + ImageUrl)
// 	})
// 	//log.Println("ASSSDAASDASD ", roomsMap)
// 	//fmt.Printf("%s", m)
// 	//var ri []RoomInfo
// 	//fmt.Println(json.Unmarshal(byteResults, &ri))
// 	return byteResults, nil
// }

// func ScrapeRooms(url string) []string {

// 	doc, err := goquery.NewDocument(url)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	var searchResults []string
// 	doc.Find("#maincontent ul.listing-results article.panel-listing-result").Each(func(i int, s *goquery.Selection) {
// 		title := s.Find("header.desktop a h1").Text()
// 		weekCost := s.Find("strong.listingPrice").First().Text()
// 		searchResults = append(searchResults, title, weekCost, "\n")
// 		fmt.Println(title + " " + weekCost)
// 	})
// 	fmt.Println(searchResults)
// 	return searchResults
// }
