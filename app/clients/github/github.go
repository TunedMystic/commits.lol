package github

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/beefsack/go-rate"
	"github.com/tunedmystic/commits.lol/app/config"
)

// Client for Github
type Client struct {
	baseURL       string
	apiKey        string
	searchLimiter *rate.RateLimiter
	maxFetch      int
	commitLength  int
}

// NewClient ...
func NewClient() *Client {
	g := Client{
		baseURL:       "https://api.github.com",
		apiKey:        config.App.GithubAPIKey,
		searchLimiter: rate.New(30-1, time.Minute),   // 30 times per minute
		maxFetch:      config.App.GithubMaxFetch,     // Max amount of items to fetch when paginating
		commitLength:  config.App.GithubCommitLength, // Max length of commit message
	}
	return &g
}

// CommitSearch ...
// The search commits endpoint works with at most 5 qualifiers
// Example:
//    https://api.github.com/search/commits?q='monkey'+author-date:2020-01-01..2020-01-13+sort:author-date-asc&page=1
//
func (g *Client) CommitSearch(options CommitSearchOptions) (*CommitSearchResponse, error) {
	// Check the rate limit, and block until the rate limit has lifted.
	g.searchLimiter.Wait()

	if options.Empty() {
		return nil, errors.New("no search options provided")
	}

	// Build request
	url := fmt.Sprintf("%v/search/commits?%v", g.baseURL, options.Serialize())

	req, _ := http.NewRequest(http.MethodGet, url, nil)

	req.Header.Add("User-Agent", "commits.lol")
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

	// Unmarshal the JSON data.
	response := CommitSearchResponse{}

	if err = json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("not able to unmarshal response: %v", err)
	}

	return &response, nil
}

// CommitSearchPaginated ...
func (g *Client) CommitSearchPaginated(options CommitSearchOptions) ([]CommitItem, error) {
	commitItems := []CommitItem{} // stores commit objects across the fetched pages

	for {
		fmt.Printf("Fetching page %v\n", options.Page)

		// Perform search.
		response, err := g.CommitSearch(options)
		if err != nil {
			return nil, err
		}

		commitItems = append(commitItems, response.CommitItems...)
		fmt.Printf("  - Page %v, got %v items\n", options.Page, len(response.CommitItems))

		// Check if last page.
		if len(commitItems) == response.TotalCount {
			fmt.Println("  - reached last page")
			break
		}

		// Check max item threshold.
		if len(commitItems) >= g.maxFetch {
			fmt.Printf("  - reached max items limit of %v\n", g.maxFetch)
			break
		}

		options.Page++
	}

	fmt.Printf("\nFetched a total of %v results\n", len(commitItems))
	return commitItems, nil
}
