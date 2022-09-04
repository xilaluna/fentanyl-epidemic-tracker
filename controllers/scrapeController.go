package controllers

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
	"github.com/gocolly/redisstorage"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
)


func ScrapeController(c *gin.Context)  {
	godotenv.Load()
	REDISHOST := os.Getenv("REDISHOST")
	REDISPORT := os.Getenv("REDISPORT")
	REDISPASSWORD := os.Getenv("REDISPASSWORD")
	
	numberOfPages := 1

	// instantiate default collector and set random User-Agent
	collector := colly.NewCollector(
		colly.AllowedDomains("darknetlive.com"),
		colly.CacheDir("./.cache"),
	)
	extensions.RandomUserAgent(collector)

	storage := &redisstorage.Storage{
		Address:  REDISHOST + ":" + REDISPORT,
		Password: REDISPASSWORD,
		DB:       0,
		Prefix:   "articleCache",
	}

	// init article collector and set random User-Agent
	articleCollector := collector.Clone()
	extensions.RandomUserAgent(articleCollector)

	err := articleCollector.SetStorage(storage)
	if err != nil {
		panic(err)
	}

	// find articles links and vist them with articleCollector
	collector.OnHTML("main > section > article > a", func(content *colly.HTMLElement) {
		link := content.Request.AbsoluteURL(content.Attr("href"))
		articleCollector.Visit(link)
	})

	// find pagination number
	collector.OnHTML("main > nav > ul > li:last-child > a", func(content *colly.HTMLElement) {
		if numberOfPages == 1 {
			re := regexp.MustCompile("[0-9]+")
			stringNumber, _ := strconv.Atoi(re.FindAllString(content.Attr("href"), -1)[0])
			numberOfPages = stringNumber
		} else {
			return
		}
	})

	articleCollector.OnHTML("main > article", func(content *colly.HTMLElement) {
		title := content.ChildText("header > h1")
		date := content.ChildText("aside > div > time")
		link := content.Request.URL.String()
		

		// loop through all the paragraphs
		content.ForEachWithBreak("div > p", func(i int, paragraph *colly.HTMLElement) bool {
			
			// check paragraph for the word "fentanyl"
			if strings.Contains(strings.ToLower(paragraph.Text), "fentanyl") {
				fmt.Println("Found article:", title, date, link)

				// Insert article into MongoDB
				document := bson.D{{Key: "link", Value: link}, {Key: "title", Value: title}, {Key: "date", Value: date}}
				articlesCollection.InsertOne(context.Background(), document)
				// stop the loop
				return false
			}
			return true
		})
	})

	articleCollector.OnRequest(func(request *colly.Request) {
		fmt.Println("Visiting", request.URL)
	})
	
	collector.OnRequest(func(request *colly.Request) {
		fmt.Println("Visiting", request.URL)
	})

	// start scraping
	for i := 1; i <= numberOfPages; i++ {
		if i == 1 {
			collector.Visit("https://darknetlive.com/post")
		} else {
			collector.Visit("https://darknetlive.com/post/page/" + strconv.Itoa(i))
			
		}
	}


	c.JSON(200, gin.H{
		"message": "scraped",
	})
}