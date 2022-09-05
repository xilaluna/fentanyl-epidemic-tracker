# [Fentanyl Epidemic Tracker](https://fentanyl-epidemic-tracker.up.railway.app/)

[![Go Report Card](https://goreportcard.com/badge/github.com/xilaluna/Fentanyl-Epidemic-Tracker)](https://goreportcard.com/report/github.com/xilaluna/Fentanyl-Epidemic-Tracker)

![graph image](/static/graph.png)

## About

A Go scraper that is executed by a gin server endpoint that scrapes and crawls through every article on https://darknetlive.com/post/. During this process, articles relating to fentanyl are automatically selected and saved to a Mongo database.

## How it Works

1. How it works is quite simple, there is a simple go scraper built using colly which visits https://darknetlive.com/post/ and its next pages. The scraper then looks into each article and searches all the paragraphs for the word "fentanyl", if the article includes the word it will save the article name and the date it was published.
2. After the data is pulled it is then put into a json file where chart.js is able to display the amount of articles per month.

![scraper image](/static/scraper-terminal.png)

## Tech Stack

- [Colly](http://go-colly.org/)
- [Gin](https://github.com/gin-gonic/gin)
- [MongoDB](https://www.mongodb.com/)
- [Railway](https://railway.app/)
- [Chart.js](https://www.chartjs.org/)
- [Tailwind](https://tailwindcss.com/)
- [Go](https://go.dev/)

## Future Plans

- Store data in a Redis database.
- Create display to show each individual article.
- Create a way to search for specific articles.
- Speed up scraper.
