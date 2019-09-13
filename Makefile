APPNAME=hooked
VERSION=v0.0.1
SOURCE=cmd
TARGET=target

PROFILE=larwef
S3_LAMBDA_BUCKET=hooked-artifact-bucket

LAMBDA_TARGET=$(TARGET)/lambda/$(APPNAME)-$(VERSION)-lambda-deployment.zip

all: build upload

# PHONY used to mitigate conflict with dir name test
.PHONY: test
test:
	go mod tidy
	go fmt ./...
	go vet ./...
	golint ./...
	go test ./...

build: test build-lambda

build-lambda:
	GOOS=linux go build -o $(TARGET)/lambda/main $(SOURCE)/lambda/main.go
	zip -j $(LAMBDA_TARGET) $(TARGET)/lambda/main

upload: upload-lambda

upload-lambda:
	aws s3 cp $(LAMBDA_TARGET) s3://$(S3_LAMBDA_BUCKET)/$(APPNAME)-$(VERSION)-lambda-deployment.zip --profile $(PROFILE)