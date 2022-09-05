package controllers

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/xilaluna/fentanyl-epidemic-tracker/configs"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var articlesCollection *mongo.Collection = configs.DatabaseCollection(configs.GetClient(), "articles")

func GetArticles(c *gin.Context) {
	ctx := context.Background()
	var articles []bson.M
	filter := bson.M{"datapoint": true}
	
	cursor, err := articlesCollection.Find(ctx, filter)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	defer cursor.Close(ctx)
	for cursor.Next(ctx) {
		var singleArticle bson.M
		err := cursor.Decode(&singleArticle)
		if err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		articles = append(articles, singleArticle)
	}
	c.JSON(200, articles)
}