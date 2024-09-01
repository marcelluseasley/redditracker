package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/marcelluseasley/redditracker/config"
	"github.com/marcelluseasley/redditracker/internal/client"
	"github.com/patrickmn/go-cache"
)

type Server struct {
	server       *http.Server
	redditClient *client.RedditClient
	resultsCache *cache.Cache

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
		resultsCache:     cache.New(10*time.Second, 1*time.Minute),
		postDataChannel:  make(chan []client.PostData, 200),
		userCountChannel: make(chan []client.UserPostCount, 200),
	}

	r.GET("/health", server.health)
	r.GET("/", server.homeHandler)
	r.GET("/ws", server.wsHandler)
	r.GET("/data", server.resultsHandler)

	return server
}

func (s *Server) Start() error {
	go func() {
		for {
			allPosts, sortedUsers := s.redditClient.GetSubredditUsers(s.redditClient.Config)

			if len(allPosts) >= 100 {
				allPosts = allPosts[:100]
			}
			s.resultsCache.Set(s.redditClient.Config.SubReddit+"#posts", allPosts[:100], cache.NoExpiration)
			s.postDataChannel <- allPosts[:100]

			if len(sortedUsers) >= 100 {
				sortedUsers = sortedUsers[:100]
			}
			s.resultsCache.Set(s.redditClient.Config.SubReddit+"#users", sortedUsers[:100], cache.NoExpiration)
			s.userCountChannel <- sortedUsers[:100]
		}
	}()

	return s.server.ListenAndServe()
}
