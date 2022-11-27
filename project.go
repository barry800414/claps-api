package main

import (
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Project struct {
	Id            string `form:"id" json:"id" xml:"id"`
	Name          string `form:"name" json:"name" xml:"name"`
	Website       string `form:"website" json:"website" xml:"website"`
	MaxClapsCount int    `form:"maxClapsCount" json:"maxClapsCount" xml:"maxClapsCount"`
	UserId        string `form:"userId" json:"userId" xml:"userId"`
}

func CreateProject(c *gin.Context) {
	var project Project
	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	project.Id = uuid.New().String()[24:]
	// FIXME
	project.UserId = uuid.New().String()
	item, err := dynamodbattribute.MarshalMap(project)
	if err != nil {
		log.Fatalf("Got error marshalling map: %s", err)
	}
	input := &dynamodb.PutItemInput{
		Item:                item,
		TableName:           aws.String("project"),
		ConditionExpression: aws.String("attribute_not_exists(id)"),
	}
	value, exists := c.Get("db")
	db := value.(*dynamodb.DynamoDB)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Unable to connect DB"})
	}
	_, err = db.PutItem(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeConditionalCheckFailedException:
				c.JSON(http.StatusForbidden, gin.H{"message": "project exists"})
				return
			default:
				log.Fatalf("Got error calling PutItem: %s", aerr.Error())
				c.JSON(http.StatusInternalServerError, gin.H{"message": aerr.Error()})
				return
			}
		}
	}
	c.JSON(http.StatusOK, project)
}

func ReadProject(c *gin.Context) {
	projectId := c.Param("pid")
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
	} else {
		project := Project{}
		dynamodbattribute.UnmarshalMap(result.Item, &project)
		c.JSON(http.StatusOK, project)
		return
	}
}
