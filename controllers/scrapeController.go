package controllers

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
	"github.com/xilaluna/fentanyl-epidemic-tracker/configs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var websiteCollection *mongo.Collection = configs.DatabaseCollection(configs.GetClient(), "websites")


func ScrapeController(c *gin.Context)  {
	var num int
	var onPageOne bool

	

	var resultNumber bson.M
	numberPageFilter := bson.M{"website": "darknetlive.com"}
	err := websiteCollection.FindOne(context.Background(), numberPageFilter).Decode(&resultNumber)
	if err != nil {
		num = 1
	} else {
		// convert bson to int
		num = int(resultNumber["paginationNum"].(int32))
	}


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

		// Check if link is already in database
		var result bson.M
		filter := bson.M{"link": link}
		err := articlesCollection.FindOne(context.Background(), filter).Decode(&result)
		if err == nil {
			return
		}

		articleCollector.Visit(link)
	})

	// find pagination number
	collector.OnHTML("main > nav > ul > li:last-child > a", func(content *colly.HTMLElement) {
		if onPageOne {
			re := regexp.MustCompile("[0-9]+")
			stringNumber, _ := strconv.Atoi(re.FindAllString(content.Attr("href"), -1)[0])
			if stringNumber > num {
				// Update pagination number
				update := bson.M{"$set": bson.M{"website": "darknetlive.com", "paginationNum": stringNumber}}
				opts := options.Update().SetUpsert(true)
				_, err := websiteCollection.UpdateOne(context.Background(), numberPageFilter, update, opts)
				if err != nil {
					fmt.Println(err)
				}
				onPageOne = false
				num = stringNumber
			} else {
				onPageOne = false
				num = 1
				return
			}
		} else {
			return
		}
	})

	articleCollector.OnHTML("main > article", func(content *colly.HTMLElement) {
		title := content.ChildText("header > h1")
		date := content.ChildText("aside > div > time")
		link := content.Request.URL.String()
		found := false
		

		// loop through all the paragraphs
		content.ForEachWithBreak("div > p", func(i int, paragraph *colly.HTMLElement) bool {
			
			// check paragraph for the word "fentanyl"
			if strings.Contains(strings.ToLower(paragraph.Text), "fentanyl") {
				fmt.Println("Found article:", title, date, link)

				// Insert article into MongoDB
				document := bson.D{{Key: "link", Value: link}, {Key: "title", Value: title}, {Key: "date", Value: date}, {Key: "datapoint", Value: true}}
				articlesCollection.InsertOne(context.Background(), document)
				found = true
				// stop the loop
				return false
			}
			return true
		})

		if !found {
			document := bson.D{{Key: "link", Value: link}, {Key: "title", Value: title}, {Key: "date", Value: date}, {Key: "datapoint", Value: false}}
			articlesCollection.InsertOne(context.Background(), document)
		}
	})

	articleCollector.OnRequest(func(request *colly.Request) {
		fmt.Println("Visiting", request.URL)
	})
	
	collector.OnRequest(func(request *colly.Request) {
		fmt.Println("Visiting", request.URL)
	})

	// start scraping
	for i := 1; i <= num; i++ {
		if i == 1 {
			onPageOne = true
			collector.Visit("https://darknetlive.com/post")
		} else {
			collector.Visit("https://darknetlive.com/post/page/" + strconv.Itoa(i))
			
		}
	}


	c.JSON(200, gin.H{
		"message": "scraped",
	})
}