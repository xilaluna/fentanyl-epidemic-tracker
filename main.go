package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)


func main() {
	// Load .env file
	godotenv.Load()

	// Create a gin router
	router := gin.Default()

	// Connect to MongoDB
	client, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("MONGO_URL")))
	if err != nil {
		log.Fatal(err)
	}

	// Connect context
	ctx := context.Background()
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// Disconnect connection when function ends
	defer client.Disconnect(ctx)
	
	// Ping the primary
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	// Get a handle for your collection
	articlesCollection := client.Database("fentanyl-epidemic-data").Collection("articles")

	router.StaticFile("/", "./static/index.html")

	router.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.GET("/data", func(c *gin.Context) {
		articles, err := articlesCollection.Find(ctx, bson.M{})
		if err != nil {
			fmt.Println(err)
		}
		
		c.JSON(200, articles)
	})

	router.GET("/scrape", func(c *gin.Context) {
		numberOfPages := 1

		// instantiate default collector and set random User-Agent
		collector := colly.NewCollector(
			colly.AllowedDomains("darknetlive.com"),
			colly.CacheDir("./.cache"),
		)
		extensions.RandomUserAgent(collector)

		// init article collector and set random User-Agent
		articleCollector := collector.Clone()
		extensions.RandomUserAgent(articleCollector)

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
			
			// loop through all the paragraphs
			content.ForEachWithBreak("div > p", func(i int, paragraph *colly.HTMLElement) bool {
				
				// check paragraph for the word "fentanyl"
				if strings.Contains(strings.ToLower(paragraph.Text), "fentanyl") {
					title := content.ChildText("header > h1")
					date := content.ChildText("aside > div > time")
					link := content.Request.URL.String()

					fmt.Println("Found article:", title, date, link)

					// Insert article into MongoDB
					document := bson.D{{Key: "title", Value: title}, {Key: "date", Value: date}}
					articlesCollection.InsertOne(ctx, document)

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
	})

	router.Run()
}