
NOCHECK=-o StrictHostKeyCheck=no
BUILDFLAGS=-buildvcs=false

.PHONY: tidy
tidy:
	go mod tidy

.PHONY: upload
upload:
	go build $(BUILDFLAGS) -o bin/upload-service ./cmd/upload

.PHONY: manage
manage:
	go build $(BUILDFLAGS) -o bin/manage-service ./cmd/manage

.PHONY: serve 
serve:
	go build $(BUILDFLAGS) -v -o bin/server-service ./cmd/serve 

.PHONY: build
build: upload manage serve

