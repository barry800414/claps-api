package main

import (
	"net/http"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/gin-gonic/gin"
)

func connectDynamo() *dynamodb.DynamoDB {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := dynamodb.New(sess)
	return svc
}

func attachDB(db *dynamodb.DynamoDB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("db", db)
	}
}

func main() {
	db := connectDynamo()
	router := gin.Default()
	router.Use(attachDB(db))

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "ok",
		})
	})

	router.POST("/project", CreateProject)
	router.GET("/project/:pid", ReadProject)
	router.POST("/project/:pid/item/:iid/claps", CreateOrUpdateClaps)

	router.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
