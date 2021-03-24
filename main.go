package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"time"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
)

func main() {
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

		c.String(http.StatusOK, month)
		c.String(http.StatusOK, "\n")
		c.String(http.StatusOK, day)
		c.String(http.StatusOK, "\n")

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

		c.String(http.StatusOK, key_str)
		c.String(http.StatusOK, "\n\n")

		available  := m[key_str]

		c.String(http.StatusOK, available)
		c.String(http.StatusOK, "\n\n")

		//test of map
		/*date := "Wednesday, March 24, 2021"
		time := "1:45 PM - 2:45 PM"
		k1 := date + time

		t := m[k1]*/

		//c.String(http.StatusOK, t)
		//c.String(http.StatusOK, "\n\n")

		c.String(http.StatusOK, "Rockne Memorial Building Slots:\n\n")

		//Display all available spots 
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
		defer resp.Body.Close()
	})

	router.Run(":" + port)
}
