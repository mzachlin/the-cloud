package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"
	"strconv"
	"net/smtp"
	//"gopkg.in/mail.v2"
	//"fmt"
	//"sync"

	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
)

func sleeper(c *gin.Context) {
	available := ""
	for {
		log.Print("hello my name is ", string(c.Query("month")))
		log.Print("and the time is ", time.Now().String())
		available = check_available(c, false)
		log.Print(available)
		time.Sleep(30 * time.Second)
	}

	//send an email
	//from := "gymspotsND@gmail.com"
	//pass := "<a1b2c3d4!>"

	// Choose auth method and set it up
	auth := smtp.PlainAuth("", "gymspotsnd@gmail.com", "a1b2c3d4!", "smtp.gmail.com")

	// Here we do it all: connect to our server, set up a message and send it
	to := []string{"clink4@nd.edu"}
	msg := []byte("To: clink4@nd.edu\r\n" +
		"Subject: Why are you not using Mailtrap yet?\r\n" +
		"\r\n" +
		"Here’s the space for our great sales pitch\r\n")
	err := smtp.SendMail("smtp.gmail.com:587", auth, "gymspotsnd@gmail.com", to, msg)
	log.Print(err)

	if err != nil {
		log.Print(err)
		log.Fatal(err)
	}

}

func check_available(c *gin.Context, do_print bool) string {

	// Handle query string parameters
	month := c.Query("month")
	day := c.Query("day")

	//HTTP Request
	resp, err := http.Get("https://recregister.nd.edu/Program/GetProgramDetails?courseId=4c286489-76bf-47ab-bac2-728e84d3fc13&semesterId=4b3cccac-5940-40b2-ac08-c201ffe58d85")
	if err != nil {
		// handle error
		log.Printf("looks like there's an error :(((( try not having errors!!")
	}

	body, err := ioutil.ReadAll(resp.Body)
	sb := string(body)

	// find matches
	re := regexp.MustCompile(`([0-9]+) spot\(s\)|No Spots Available`)
	matches := re.FindAll([]byte(sb), -1)
	re2 := regexp.MustCompile(`([0-9])+:([0-9])([0-9]) [A|P]M - ([0-9])+:([0-9])([0-9]) [A|P]M`) //9:30 PM - 10:30 PM
	matches2 := re2.FindAll([]byte(sb), -1)
	re3 := regexp.MustCompile(`([A-Z]|[a-z])+day, ([A-Z]|[a-z])+ [0-9]+, 2021`)
	matches3 := re3.FindAll([]byte(sb), -1)

	//create hash table of all available spots
	var m map[string]string
	m = make(map[string]string)

	for i, s := range matches {
		key := string(matches3[i]) + string(matches2[i])
		m[key] = string(s)
	}

	const layout = "01-02-2006"
	d := string(day)
	if len(d) == 1 {
		d = "0" + d
	}
	mo := string(month)
	if len(mo) == 1 {
		mo = "0" + mo
	}
	date_s:= mo + "-" + d + "-" + "2021"
	t, _ := time.Parse(layout, date_s)
	day_str := t.Weekday().String()

	m_int, _ := strconv.Atoi(month)
	key_str := day_str + ", " + time.Month(m_int).String() + " " + string(day) + ", 2021" + "1:30 PM - 2:30 PM"
	available  := m[key_str]

	if do_print {

		c.String(http.StatusOK, key_str)
		c.String(http.StatusOK, "\n\n")

		c.String(http.StatusOK, available)
		c.String(http.StatusOK, "\n\n")


		c.String(http.StatusOK, "Rockne Memorial Building Slots:\n\n")

		//Display all spots
		for i, s := range matches {
			s := string(s)
			if s == "No Spots Available" {
				c.String(http.StatusOK, "0 spot(s)")
			} else {
				c.String(http.StatusOK, s)
			}
			c.String(http.StatusOK, "\t\t")
			c.String(http.StatusOK, string(matches2[i]))
			c.String(http.StatusOK, "\t\t")
			c.String(http.StatusOK, string(matches3[i]))
			c.String(http.StatusOK, "\n")
		}
	}
	defer resp.Body.Close()
	return available
}

func main() {
	//var wg sync.WaitGroup

	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
		log.Print("Welcome to the home page.")
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})

	router.GET("/get", func(c *gin.Context) {
		//wg.Add(1)
		check_available(c, true)
		go sleeper(c)
	})

	router.Run(":" + port)
	//wg.Wait()
}
