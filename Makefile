
NOCHECK=-o StrictHostKeyCheck=no

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: upload
upload:
	go build -o bin/upload-service ./cmd/upload

.PHONY: manage
manage:
	go build -o bin/manage-service ./cmd/manage

.PHONY: serve 
serve:
	go build -v -o bin/server-service ./cmd/serve 

.PHONY: build
build: upload manage serve

