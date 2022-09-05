# [Fentanyl Epidemic Tracker](https://fentanyl-epidemic-tracker.up.railway.app/)

[![Go Report Card](https://goreportcard.com/badge/github.com/xilaluna/Fentanyl-Epidemic-Tracker)](https://goreportcard.com/report/github.com/xilaluna/Fentanyl-Epidemic-Tracker)

![graph image](/static/graph.png)

## About

A Go scraper that is executed by a gin server endpoint that scrapes and crawls through every article on https://darknetlive.com/post/. During this process, articles relating to fentanyl are automatically selected and saved to a Mongo database.

## Usage

You are able to call the gin server endpoint to scrape the website and fetch the articles from the database. The gin server is hosted on Railway and can be accessed at https://fentanyl-epidemic-tracker.xilaluna.com/.

### API Endpoints:

- `/scrape` - Scrapes the website and saves the articles to the database.
- `/articles` - Fetches the articles from the database.
- `/ping` - Pings the gin server to check if it is up.

## How it Works

1. First the gin server is started and the database is connected.
2. The `/scrape` endpoint is called and checks pagination to see if a new page has been added.
3. If a new page has been added, the articles are scraped and saved to the database if they contain the keyword "fentanyl".
4. The `/articles` endpoint is called and the articles are fetched from the database.
5. The articles are displayed on the website.

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
- Update the asthetics of the static website.
- Create display to show each individual article.
- Create a way to search for specific articles.
- Speed up scraper.
