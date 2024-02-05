package main

import (
	"flag"

	"github.com/marcelluseasley/redditracker/api"
	"github.com/marcelluseasley/redditracker/config"
	"github.com/marcelluseasley/redditracker/internal/client"
	log "github.com/sirupsen/logrus"
)

func main() {

	conf := parseFlagsAndLoadConfig()
	token := getToken(conf)
	updateConfigWithToken(conf, token)
	startServer(&conf.Port, conf)

}

func parseFlagsAndLoadConfig() *config.Config {
	port := flag.Int("port", 8080, "server port")
	subReddit := flag.String("subreddit", "TaylorSwift", "subreddit to track")

	flag.Parse()

	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	conf.SubReddit = *subReddit
	conf.Port = *port

	return conf
}

func getToken(conf *config.Config) *client.RedditToken {
	tokenClient := client.NewTokenClient()
	token, err := tokenClient.GetToken(conf)
	if err != nil {
		log.Fatal(err)
	}
	return token
}

func updateConfigWithToken(conf *config.Config, token *client.RedditToken) {
	conf.InitialToken = token.AccessToken
	conf.ExpiresIn = token.ExpiresIn
}

func startServer(port *int, conf *config.Config) {
	server := api.NewServer(*port, conf)
	log.Println("Starting server on port", *port)
	log.Fatal(server.Start())
}
