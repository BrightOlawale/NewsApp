package news

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

// Client : represents a client that can make requests to the News API
type Client struct {
	http     *http.Client // http client
	key      string       // api key
	PageSize int          // page number
}

// Article : represents an array of articles retrieved from the News API
type Article struct {
	Source struct {
		ID   any    `json:"id"`
		Name string `json:"name"`
	} `json:"source"`
	Author      string    `json:"author"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	URL         string    `json:"url"`
	URLToImage  string    `json:"urlToImage"`
	PublishedAt time.Time `json:"publishedAt"`
	Content     string    `json:"content"`
}

// Results : represents a single response of retrieved from the News API
type Results struct {
	Status       string    `json:"status"`
	TotalResults int       `json:"totalResults"`
	Article      []Article `json:"articles"`
}

// NewClient : returns a new client instance used to make requests to the News API
func NewClient(httpClient *http.Client, key string, pageSize int) *Client {
	if pageSize > 100 {
		pageSize = 100
	}

	return &Client{http: httpClient, key: key, PageSize: pageSize}
}

func (c *Client) FetchEverything(query, page string) (*Results, error) {
	endpoint := fmt.Sprintf("https://newsapi.org/endpoints/everything?q=%s&pageSize=%d&page=%s&apiKey=%s&sortBy=publishedAt&language=en", url.QueryEscape(query), c.PageSize, page, c.key)

	resp, err := c.http.Get(endpoint)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusOK {
		return nil, fmt.Errorf(string(body))
	}

	res := &Results{}

	return res, json.Unmarshal(body, res)
}
