# Scripts to handle catena build and installation
# Shell to use with Make
SHELL := /bin/bash

# Build Environment
PACKAGE = catena
PBPKG = $(CURDIR)/api

# Commands
GOCMD = go
GODOC = godoc
PROTOC = protoc
GORUN = $(GOCMD) run
GOGET = $(GOCMD) get
GOTEST = $(GOCMD) test
GOINSTALL = $(GOCMD) install
GOCLEAN = $(GOCMD) clean
GOGENERATE = $(GOCMD) generate

# Output Helpers
BM  = $(shell printf "\033[34;1m●\033[0m")
GM = $(shell printf "\033[32;1m●\033[0m")
RM = $(shell printf "\033[31;1m●\033[0m")


# Export targets not associated with files.
.PHONY: all install catena test citest clean doc

# Ensure dependencies are installed, run tests and compile
all: test install

# Build the various binaries and sources
install: generate catena

# Build and install the catena command in the $GOBIN directory
catena:
	$(info $(GM) compiling catena executable with go install …)
	@ $(GOINSTALL) ./cmd/catena

# Run go generate to build protocol buffers and other files
generate:
	$(info $(BM) running go generate …)
	@ $(GOGENERATE) ./...

# Target for simple testing on the command line
test:
	$(info $(BM) running simple local tests …)
	@ $(GOTEST) -v ./...

# Target for testing in continuous integration
citest:
	$(info $(BM) running CI tests with randomization and race …)
	$(GOTEST) -bench=. -v --cover -coverprofile=coverage.txt -covermode=atomic --race ./...

# Run Godoc server and open browser to the documentation
doc:
	$(info $(BM) running go documentation server at http://localhost:6060)
	$(info $(BM) type CTRL+C to exit the server)
	@ open http://localhost:6060/pkg/github.com/bbengfort/catena/
	@ $(GODOC) --http=:6060

# Clean build files
clean:
	$(info $(RM) cleaning up build …)
	@ $(GOCLEAN)
	@ find . -name "*.coverprofile" -print0 | xargs -0 rm -rf
	@ find . -name "coverage.text" -print0 | xargs -0 rm -rf
