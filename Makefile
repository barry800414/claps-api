
deploy-prd: 
GOOS=linux GOARCH=arm64 go build -o bootstrap
zip bootstrap.zip bootstrap
