package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
	"github.com/joho/godotenv"
)
type Article struct {
	Title string `json:"title"`
	Date string `json:"date"`
}

var ctx = context.Background()


func main() {
	// Load .env file
	godotenv.Load()
	
	// Connect to Redis
	url := os.Getenv("REDIS_URL")
	opt, err := redis.ParseURL(url)
	if err != nil {
		panic(err)
	}
	rdb := redis.NewClient(opt)

	// Create a gin router
	router := gin.Default()

	router.StaticFile("/", "./static/index.html")

	router.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.GET("/articles", func(c *gin.Context) {
		// Get the articles from Redis
		articles, err := rdb.Get(ctx, "articles").Result()
		if err != nil {
			panic(err)
		}
		c.JSON(200, gin.H{
			"articles": articles,
		})
	})

	router.GET("/scrape", func(c *gin.Context) {
		articles := []Article{}
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

					fmt.Println("Found article:", title, date)

					// add article to articles slice
					newArticle := Article{}
					newArticle.Title = title
					newArticle.Date = date
					
					articles = append(articles, newArticle)

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

		// convert articles to json
		json, err := json.Marshal(articles)
		if err != nil {
			log.Println(err)
		}

		err = rdb.Set(ctx, "articles", json, 0).Err()
		if err != nil {
			fmt.Println(err)
		}


		c.JSON(200, gin.H{
			"message": "scraped",
		})
	})

	router.Run()
}
