#!make

CTR_REGISTRY    ?= allenlsy
CTR_TAG         ?= latest
GIT_SHA=$$(git rev-parse HEAD)

build:
	go build -v -o ./bin/client-app ./cmd/client-app
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -o ./bin/client-app-linux ./cmd/client-app

docker-build: build
	docker build . -t $(CTR_REGISTRY)/simple-request-app:$(CTR_TAG)
	docker build . -t $(CTR_REGISTRY)/simple-request-app:$(GIT_SHA)

docker-push: docker-build
	docker push $(CTR_REGISTRY)/simple-request-app:$(CTR_TAG)
	docker push $(CTR_REGISTRY)/simple-request-app:$(GIT_SHA)