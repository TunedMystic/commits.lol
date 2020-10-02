package github

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

// Client for Github
type Client struct {
	apiKey  string
	baseURL string
}

// NewClient ...
func NewClient(apiKey string) *Client {
	g := Client{
		apiKey:  apiKey,
		baseURL: "https://api.github.com",
	}
	return &g
}

// SetBaseURL ...
func (g *Client) SetBaseURL(url string) {
	g.baseURL = url
}

// CommitSearch ...
// The search commits endpoint works with at most 5 qualifiers
// Example:
//    https://api.github.com/search/commits?q='monkey'+author-date:2020-01-01..2020-01-13+sort:author-date-asc&page=1
//
func (g *Client) CommitSearch(options CommitSearchOptions) (*CommitSearchResponse, error) {
	if options.Empty() {
		return nil, errors.New("must provide commit search options")
	}

	// Build request
	url := fmt.Sprintf("%v/search/commits?%v", g.baseURL, options.Serialize())
	fmt.Println(url)
	// url := string(0x7f)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	fmt.Println(err)
	if err != nil {
		return nil, fmt.Errorf("error building the request: %v", err)
	}

	req.Header.Add("Accept", "application/vnd.github.cloak-preview+json")
	req.Header.Add("Authorization", "token "+g.apiKey)

	// Make request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}

	// Read the response body.
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading the response body: %v", err)
	}
	res.Body.Close()

	fmt.Printf(string(data))

	if res.StatusCode != 200 {
		return nil, fmt.Errorf("got response %v", res.StatusCode)
	}

	// Unmarshal the JSON data.
	response := CommitSearchResponse{}

	if err = json.Unmarshal(data, &response); err != nil {
		fmt.Println("Error unmarshalling the json data")
		return nil, err
	}

	return &response, nil
}
