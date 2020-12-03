#
# Makefile to perform "live code reloading" after changes to .go files.
#
# n.b. you must install fswatch (OS X: `brew install fswatch`)
#
# To start live reloading run the following command:
# $ make serve
#

mb_version := $(shell cat ../VERSION)
mb_count := $(shell git rev-list HEAD --count)
mb_hash := $(shell git rev-parse --short HEAD)

# binary name to kill/restart
PROG = showman
 
# targets not associated with files
.PHONY: default build
 
# clean up
clean:
	go clean

# run formatting tool and build
build: clean
	go build fmt
	go build -ldflags "-s -w -X main.Version=$(mb_version)-$(mb_count).$(mb_hash)" -v -o ${PROG}

release: clean
	go build fmt
	go build -ldflags "-s -w -X main.Version=$(mb_version)-$(mb_count).$(mb_hash)" -v -o ${PROG}
