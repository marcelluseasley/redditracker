PORT ?= 8080
SUBREDDIT ?= TaylorSwift

build:
	go build -o cmd/redditracker/redditracker cmd/redditracker/main.go

run: build
	cd cmd/redditracker && GIN_MODE=release ./redditracker -port=$(PORT) -subreddit=$(SUBREDDIT)

clean:
	rm cmd/redditracker/redditracker