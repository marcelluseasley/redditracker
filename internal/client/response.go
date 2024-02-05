package client

type RedditPostResponse struct {
	Data RedditData `json:"data"`
}

type RedditData struct {
	After    string       `json:"after"`
	Children []RedditPost `json:"children"`
	Before   *string      `json:"before"`
}

type RedditPost struct {
	Data PostData `json:"data"`
}

type PostData struct {
	Subreddit  string  `json:"subreddit"`
	Author     string  `json:"author"`
	Title      string  `json:"title"`
	Ups        int     `json:"ups"`
	Created    float64 `json:"created"`
	CreatedUTC float64 `json:"created_utc"`
	ID         string  `json:"id"`
}
