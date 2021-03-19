package main

import (
  "bytes"
	"log"
	"net/http"
	"os"
  "strconv"
	"io/ioutil"
	"regexp"

	"github.com/gin-gonic/gin"
	_ "github.com/heroku/x/hmetrics/onload"
  "github.com/russross/blackfriday"
)

func repeatHandler(r int) gin.HandlerFunc {
  return func(c *gin.Context) {
    var buffer bytes.Buffer
    for i := 0; i < r; i++ {
      buffer.WriteString("Hello from Go!\n")
    }
    c.String(http.StatusOK, buffer.String())
  }
}


func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

  tStr := os.Getenv("REPEAT")
  repeat, err := strconv.Atoi(tStr)
  if err != nil {
    log.Printf("Error converting $REPEAT to an int: %q - Using default\n", err)
  }

	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
		log.Print("Welcome to the home page.")
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})

  router.GET("/mark", func(c *gin.Context) {
		c.String(http.StatusOK, "hi there")
  })

	router.GET("/two", func(c *gin.Context) {
		c.String(http.StatusOK, string(blackfriday.Run([]byte("Well hello there"))))
	})

  router.GET("/repeat", repeatHandler(repeat))

	router.GET("/get", func(c *gin.Context) {
		resp, err := http.Get("https://recregister.nd.edu/Program/GetProgramDetails?courseId=4c286489-76bf-47ab-bac2-728e84d3fc13&semesterId=4b3cccac-5940-40b2-ac08-c201ffe58d85")
		if err != nil {
			// handle error
			log.Printf("looks like there's an error :(((( lmao")
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

	c.String(http.StatusOK, "Rockne Memorial Building Slots:\n\n")
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
