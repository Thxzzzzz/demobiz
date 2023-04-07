IMG ?= 379809513/demo-biz:latest

run:
	go run main.go

docker-build:  ## Build docker image with the manager.
	docker build -t ${IMG} .

docker-push: ## Push docker image with the manager.
	docker push ${IMG}

docker-run:
	docker run --name demo-sidecar ${IMG} --port 8080:80