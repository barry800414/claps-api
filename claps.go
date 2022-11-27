package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type PostClapBody struct {
	Claps int `form:"claps" json: "claps" xml:"claps" binding:"required"`
}

func CreateOrUpdateClaps(c *gin.Context) {
	projectId := c.Param("pid")
	itemId := c.Param("pid")
	projectItemId := fmt.Sprintf("%s#%s", projectId, itemId)
	cuuid, err := c.Cookie("clapsUuid")
	var postClapBody PostClapBody
	if err := c.ShouldBindJSON(&postClapBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if cuuid == "" {
		cuuid = uuid.New().String()
		c.SetCookie("clapsUUid", cuuid, 31536000, "/", "getclaps.com", false, true)
	}

	value, exists := c.Get("db")
	db := value.(*dynamodb.DynamoDB)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Unable to connect DB"})
	}
	result, err := db.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("project"),
		Key: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(projectId),
			},
		},
	})
	if err != nil {
		log.Fatalf("Got error calling GetItem: %s", err)
	}

	if result.Item == nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "project not found"})
		return
	}
	project := Project{}
	dynamodbattribute.UnmarshalMap(result.Item, &project)
	fmt.Printf("project: %v\n", project)
	// maxClapsCount := 1
	// if project.MaxClapsCount != 0 {
	// 	maxClapsCount = project.MaxClapsCount
	// }

	result, err = db.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("userClap"),
		Key: map[string]*dynamodb.AttributeValue{
			"projectItemId": {
				S: aws.String(projectItemId),
			},
			"cuuid": {
				S: aws.String(cuuid),
			},
		},
	})
	if result.Item == nil {
		// create new item in userClap table
		// create or update value in clap table
	} else {
		// update item in userClap table
		// create or update value in clap table
	}
	fmt.Printf("%v %v", result, result.Item)

}
