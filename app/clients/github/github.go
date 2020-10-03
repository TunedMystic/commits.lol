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

// CommitSearch ...
// The search commits endpoint works with at most 5 qualifiers
// Example:
//    https://api.github.com/search/commits?q='monkey'+author-date:2020-01-01..2020-01-13+sort:author-date-asc&page=1
//
func (g *Client) CommitSearch(options CommitSearchOptions) (*CommitSearchResponse, error) {
	if options.Empty() {
		return nil, errors.New("no search options provided")
	}

	// Build request
	url := fmt.Sprintf("%v/search/commits?%v&per_page=1", g.baseURL, options.Serialize())

	req, _ := http.NewRequest(http.MethodGet, url, nil)

	req.Header.Add("Accept", "application/vnd.github.cloak-preview+json")
	req.Header.Add("Authorization", "token "+g.apiKey)

	// Make request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}

	// Read the response body.
	data, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, NewAPIError(url, data, res.StatusCode)
	}

	fmt.Println(string(data))

	// Unmarshal the JSON data.
	response := CommitSearchResponse{}

	if err = json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("not able to unmarshal response: %v", err)
	}

	return &response, nil
}
