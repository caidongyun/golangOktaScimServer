# This makefile was borrowed from pmcdowell@okta.com, and I appreciate you using it.
# Feel free to make it better !

GOPATH=$(shell pwd)
SHELL := /bin/bash
PATH := bin:$(PATH)

# This is the library that provided the LDAP capability
setup:
	@GOPATH=$(GOPATH) go get -u github.com/gorilla/mux
	#@GOPATH=$(GOPATH) go get "github.com/vjeantet/goldap/message"

build:
	make setup
	env GOOS=linux GOARCH=386 go build  -o scimServer.linux main.go structs.go
	env GOOS=darwin GOARCH=386 go build -o scimServer.macos main.go structs.go
	env GOOS=windows GOARCH=386 go build -o scimServer.exe main.go structs.go

	chmod +x okta2anything.linux
	chmod +x okta2anything.macos
	make push

buildmacos:
	env GOOS=darwin GOARCH=386 go build -o okta2anything.macos okta2anything.go
	make push
clean:
	rm okta2anything.linux
	rm okta2anything.exe
	rm okta2anything.macos

push:
	git add *
	git commit -m "push"
	git push origin master
    
