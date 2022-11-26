package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Project struct {
	Id            string `form:"id" json:"id" xml:"id"`
	Name          string `form:"name" json:"name" xml:"name"`
	Website       string `form:"website" json:"website" xml:"website"`
	MaxClapsCount int    `form:"max_claps_count" json:"max_claps_count" xml:"max_claps_count"`
}

func connectDynamo() *dynamodb.DynamoDB {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := dynamodb.New(sess)
	return svc
}

func main() {
	db := connectDynamo()
	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
		})
	})

	router.POST("/project", func(c *gin.Context) {
		var project Project
		if err := c.ShouldBindJSON(&project); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		project.Id = uuid.New().String()[24:]
		item, err := dynamodbattribute.MarshalMap(project)
		if err != nil {
			log.Fatalf("Got error marshalling map: %s", err)
		}
		input := &dynamodb.PutItemInput{
			Item:      item,
			TableName: aws.String("project"),
		}
		x, err := db.PutItem(input)
		if err != nil {
			log.Fatalf("Got error calling PutItem: %s", err)
		}
		fmt.Printf("%T %v", x, x)
		c.JSON(http.StatusOK, gin.H{"name": project.Name, "website": project.Website, "max_claps_count": project.MaxClapsCount})
	})

	router.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
