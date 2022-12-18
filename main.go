package main

import (
	"context"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	ginadapter "github.com/awslabs/aws-lambda-go-api-proxy/gin"
	"github.com/gin-gonic/gin"
)

var ginLambda *ginadapter.GinLambda

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

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return ginLambda.ProxyWithContext(ctx, request)
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

	env := os.Getenv("GIN_MODE")
	if env == "release" {
		ginLambda = ginadapter.New(router)
		lambda.Start(Handler)
	} else {
		router.Run(":8080")
	}
}
