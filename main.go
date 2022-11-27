package main

import (
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
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
			Item:                item,
			TableName:           aws.String("project"),
			ConditionExpression: aws.String("attribute_not_exists(id)"),
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
		} else {
			log.Fatalf("Got error calling PutItem: %s", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"name": project.Name, "website": project.Website, "max_claps_count": project.MaxClapsCount})
	})

	router.GET("/project/:pid", func(c *gin.Context) {
		projectId := c.Param("pid")

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
	})

	router.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
