# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOINSTALL=$(GOCMD) install

build: dep doc
	$(GOBUILD) -o server ./main.go

dep: wire
	$(shell $(GOCMD) env GOPATH)/bin/wire ./internal/wire
wire:
	GO111MODULE=on $(GOINSTALL) github.com/google/wire/cmd/wire@v0.6.0
