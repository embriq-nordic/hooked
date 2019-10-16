APPNAME=hooked
VERSION=v0.1.0
SOURCE=cmd
TARGET=target
PORT=8081

GOOS=linux
GOARCH=amd64

PROFILE=larwef
S3_LAMBDA_BUCKET=hooked-bucket

LAMBDA_TARGET=$(TARGET)/lambda/$(APPNAME)-$(VERSION)-lambda-deployment.zip

all: build upload

coverage:
	go test ./... -coverprofile=coverage.out
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out

# PHONY used to mitigate conflict with dir name test
.PHONY: test
test:
	go mod tidy
	go fmt ./...
	go vet ./...
	golint ./...
	go test ./...

integration:
	go test ./... -tags=integration

# Run locally
docker: build-server build-docker run-docker

run-docker:
	docker run -it --rm -p $(PORT):$(PORT) \
	-e port=$(PORT) \
	$(APPNAME)

# Build
build: test build-lambda build-server build-docker

build-lambda:
	GOOS=linux go build -o $(TARGET)/lambda/main $(SOURCE)/lambda/main.go
	zip -j $(LAMBDA_TARGET) $(TARGET)/lambda/main

build-server:
	GOOS=$(GOOS) go build -ldflags "-X main.version=$(VERSION)" -o target/server/app cmd/server/main.go

build-docker:
	docker build -t $(APPNAME) -f build/docker/Dockerfile .

# Upload
upload: upload-lambda

upload-lambda:
	aws s3 cp $(LAMBDA_TARGET) s3://$(S3_LAMBDA_BUCKET)/$(APPNAME)-$(VERSION)-lambda-deployment.zip --profile $(PROFILE)

clean:
	rm -rf $(TARGET)

rebuild:
	clean all

doc:
	godoc -http=":6060"
