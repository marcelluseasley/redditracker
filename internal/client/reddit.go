package client

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/marcelluseasley/redditracker/config"
)

const (
	NumOfHTTPCalls = 3
)

var wg sync.WaitGroup

type UserPostCount struct {
	Username  string
	PostCount int
}

type RedditClient struct {
	Config     *config.Config
	httpClient *http.Client
}

func NewRedditClient(conf *config.Config) *RedditClient {
	jwtTransport := NewJWTTransport(conf.InitialToken, ExpiresInToExpiry(conf.ExpiresIn), conf)

	return &RedditClient{
		Config: conf,
		httpClient: &http.Client{
			Timeout:   5 * time.Second,
			Transport: jwtTransport,
		},
	}
}

func ExpiresInToExpiry(expiresIn int) time.Time {
	return time.Now().Add(time.Duration(expiresIn) * time.Second)
}

// calls to reddit

func (c *RedditClient) GetSubredditUsers(conf *config.Config) ([]PostData, []UserPostCount) {
	allPosts := &[]PostData{}
	userMap := make(map[string][]string)

	done := make(chan struct{})
	userDataChannel := make(chan PostData)
	postDataChannel := make(chan PostData)

	wg.Add(3)
	go func() {
		defer wg.Done()
		getSubredditTop(c, conf, userDataChannel, postDataChannel)
	}()

	go func() {
		defer wg.Done()
		getSubredditUsersNew(c, conf, userDataChannel)
	}()

	go func() {
		defer wg.Done()
		getSubredditUsersHot(c, conf, userDataChannel)
	}()

	go func() {
		wg.Wait()
		close(userDataChannel)
		close(postDataChannel)
		close(done)
	}()

loop:
	for {
		select {

		case post := <-postDataChannel:
			*allPosts = append(*allPosts, post)
		case userPost := <-userDataChannel:
			if userPosts, ok := userMap[userPost.Author]; !ok {
				userMap[userPost.Author] = []string{userPost.ID}
			} else {
				if !postInUserMap(userPost.ID, userPosts) {
					userMap[userPost.Author] = append(userMap[userPost.Author], userPost.ID)
				}
			}
		case <-done:
			break loop
		}
	}

	usersSortedByPostCount := sortUserByNumPosts(userMap)

	return *allPosts, usersSortedByPostCount
}

func postInUserMap(postID string, posts []string) bool {
	for _, id := range posts {
		if id == postID {
			return true
		}
	}
	return false
}

func getSubredditTop(c *RedditClient, conf *config.Config, usersChan chan PostData, postsChan chan PostData) {
	getSubredditData(c, conf, "top.json", usersChan, postsChan)
}

func getSubredditUsersNew(c *RedditClient, conf *config.Config, usersChan chan PostData) {
	getSubredditData(c, conf, "new.json", usersChan, nil)
}

func getSubredditUsersHot(c *RedditClient, conf *config.Config, usersChan chan PostData) {
	getSubredditData(c, conf, "hot.json", usersChan, nil)
}

func getSubredditData(c *RedditClient, conf *config.Config, path string, usersChan, postsChan chan PostData) {
	// for paging
	after := ""
	for {
		req, err := http.NewRequest("GET", fmt.Sprintf("%s/r/%s/%s?limit=100&t=all&after=%s", baseURL, conf.SubReddit, path, after), nil)
		if err != nil {
			log.Println(err)
			return
		}

		req.SetBasicAuth(conf.RedditClientID, conf.RedditClientSecret)
		req.Header.Set("User-Agent", conf.UserAgent)
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		resp, err := c.httpClient.Do(req)
		if err != nil {
			log.Println(err)
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
			return
		}
		var posts RedditPostResponse
		err = json.Unmarshal(body, &posts)
		if err != nil {
			log.Println(err)
			return
		}

		for _, post := range posts.Data.Children {
			if usersChan != nil {
				usersChan <- post.Data
			}
			if postsChan != nil {
				postsChan <- post.Data
			}

		}

		// after = posts.Data.After

		if after == "" {
			break
		}
		requestDelay(resp.Header.Get("X-Ratelimit-Remaining"), resp.Header.Get("X-Ratelimit-Reset"))

	}
}

func sortUserByNumPosts(userMap map[string][]string) []UserPostCount {
	var userPostCount []UserPostCount
	for key, value := range userMap {
		userPostCount = append(userPostCount, UserPostCount{Username: key, PostCount: len(value)})
	}
	sort.Slice(userPostCount, func(i, j int) bool {
		return userPostCount[i].PostCount > userPostCount[j].PostCount
	})
	return userPostCount
}

func requestDelay(xRateLimitRemaining string, xRateLimitReset string) {
	limitRemaining, err := strconv.Atoi(xRateLimitRemaining)
	if err != nil {
		log.Println(err)
		return
	}
	limitReset, err := strconv.Atoi(xRateLimitReset)
	if err != nil {
		log.Println(err)
		return
	}
	time.Sleep(time.Duration(limitReset/(limitRemaining/NumOfHTTPCalls)) * time.Second)
}
