package main

import (
  "bytes"
	"log"
	"net/http"
	"os"
  "strconv"

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

func handleDateEntry() gin.HandlerFunc {
  return func(c *gin.Context) {
    log.Print("You have selected a month!")
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

  // TODO: get rid of this later
  router.GET("/mark", func(c *gin.Context) {
    c.String(http.StatusOK, string(blackfriday.Run([]byte("**hi!**"))))
  })

  // TODO: get rid of this later
  router.GET("/repeat", repeatHandler(repeat))

	router.Run(":" + port)
}
