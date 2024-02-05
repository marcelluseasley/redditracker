package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/marcelluseasley/redditracker/config"
	"github.com/marcelluseasley/redditracker/internal/client"
)

type Server struct {
	server       *http.Server
	redditClient *client.RedditClient

	postDataChannel  chan []client.PostData
	userCountChannel chan []client.UserPostCount
}

func NewServer(port int, conf *config.Config) *Server {

	r := gin.Default()

	hs := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	redditClient := client.NewRedditClient(conf)

	server := &Server{
		server:           hs,
		redditClient:     redditClient,
		postDataChannel:  make(chan []client.PostData),
		userCountChannel: make(chan []client.UserPostCount),
	}

	r.GET("/health", server.health)
	r.GET("/", server.homeHandler)
	r.GET("/ws", server.wsHandler)

	return server
}

func (s *Server) Start() error {
	go func() {
		for {
			allPosts, sortedUsers := s.redditClient.GetSubredditUsers(s.redditClient.Config)

			if len(allPosts) > 100 {
				s.postDataChannel <- allPosts[:100]
			}

			if len(sortedUsers) > 100 {
				s.userCountChannel <- sortedUsers[:100]
			}

		}
	}()

	return s.server.ListenAndServe()
}
