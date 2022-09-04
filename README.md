# [Fentanyl Epidemic Tracker](https://fentanyl-epidemic-tracker.up.railway.app/)

[![Go Report Card](https://goreportcard.com/badge/github.com/xilaluna/Fentanyl-Epidemic-Tracker)](https://goreportcard.com/report/github.com/xilaluna/Fentanyl-Epidemic-Tracker)

The mission for this project was to create a scraper that would pull articles pertaining to fentanyl, for the goal of mapping the trend of the fentanyl epidemic & the increase of fentanyl distrubtion through illegal sites.

![graph image](/static/graph.png)

## How it Works

1. How it works is quite simple, there is a simple go scraper built using colly which visits https://darknetlive.com/post/ and its next pages. The scraper then looks into each article and searches all the paragraphs for the word "fentanyl", if the article includes the word it will save the article name and the date it was published.
2. After the data is pulled it is then put into a json file where chart.js is able to display the amount of articles per month.

![scraper image](/static/scraper-terminal.png)

## Tech Stack

- [Go](https://go.dev/)
- [Colly](http://go-colly.org/)
- [Gin](https://github.com/gin-gonic/gin)
- [Redis](https://redis.io/)
- [MongoDB](https://www.mongodb.com/)
- [Railway](https://railway.app/)
- [Chart.js](https://www.chartjs.org/)
- JSON
- HTML

## Future Plans

- Store data in a database such as Redis.
- Store every article in the database then check if the link exists already for the collector.
- Create display to show each individual article.
- Create a way to search for specific articles.
- Speed up scraper.
