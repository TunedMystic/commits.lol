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
	"go.uber.org/zap"
)

// Client for Github
type Client struct {
	baseURL          string
	apiKey           string
	searchLimiterMin *rate.RateLimiter
	searchLimiterHr  *rate.RateLimiter
	maxFetch         int
	commitLength     int
}

// NewClient ...
func NewClient() Client {
	return Client{
		baseURL:          "https://api.github.com",
		apiKey:           config.App.GithubAPIKey,
		searchLimiterMin: rate.New(30, time.Second*70),   // 30 times per 70 seconds
		searchLimiterHr:  rate.New(5000, time.Minute*70), // 5000 times per 70 minutes
		maxFetch:         config.App.GithubMaxFetch,      // Max amount of items to fetch when paginating
		commitLength:     config.App.GithubCommitLength,  // Max length of commit message
	}
}

// RateLimits checks the rate limit for the configured API Key.
func (g *Client) RateLimits() (RateLimitResponse, error) {
	// Build request
	url := fmt.Sprintf("%v/rate_limit", g.baseURL)
	req, _ := http.NewRequest(http.MethodGet, url, nil)

	req.Header.Add("User-Agent", "commits.lol")
	req.Header.Add("Accept", "application/vnd.github.v3+json")
	req.Header.Add("Authorization", "token "+g.apiKey)

	var response RateLimitResponse

	// Make request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return response, fmt.Errorf("error making request: %v", err)
	}

	// Read the response body.
	data, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return response, NewAPIError(url, data, res.StatusCode)
	}

	// Unmarshal the JSON data.
	if err = json.Unmarshal(data, &response); err != nil {
		return response, fmt.Errorf("not able to unmarshal response: %v", err)
	}

	return response, nil
}

// CommitSearch ...
// The search commits endpoint works with at most 5 qualifiers
// Example:
//    https://api.github.com/search/commits?q='monkey'+author-date:2020-01-01..2020-01-13+sort:author-date-asc&page=1
//
func (g *Client) CommitSearch(options CommitSearchOptions) (CommitSearchResponse, error) {
	response := CommitSearchResponse{}

	// Check the rate limit, and block until the rate limit has lifted.
	g.searchLimiterMin.Wait()

	if options.IsEmpty() {
		return response, errors.New("no search options provided")
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
		return response, fmt.Errorf("error making request: %v", err)
	}

	// Read the response body.
	data, _ := ioutil.ReadAll(res.Body)
	res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return response, NewAPIError(url, data, res.StatusCode)
	}

	// Unmarshal the JSON data.
	if err = json.Unmarshal(data, &response); err != nil {
		return response, fmt.Errorf("not able to unmarshal response: %v", err)
	}

	return response, nil
}

// CommitSearchPaginated ...
func (g *Client) CommitSearchPaginated(options CommitSearchOptions) ([]CommitItem, error) {
	commitItems := make([]CommitItem, 0, 30) // stores commit objects across the fetched pages

	for {
		zap.S().Infof("  Query [%s], fetching Page %d", options.QueryText, options.Page)

		// Perform search.
		response, err := g.CommitSearch(options)
		if err != nil {
			return nil, err
		}

		commitItems = append(commitItems, response.CommitItems...)
		zap.S().Debugf("    - Query [%s], Page %d, got %d items", options.QueryText, options.Page, len(response.CommitItems))

		// Check if last page.
		if len(commitItems) == response.TotalCount {
			zap.S().Debugf("    - Query [%s], reached last Page %d", options.QueryText, options.Page)
			break
		}

		// Check max item threshold.
		if len(commitItems) >= g.maxFetch {
			zap.S().Debugf("    - Query [%s], reached items limit of %d", options.QueryText, g.maxFetch)
			break
		}

		options.Page++
	}

	zap.S().Infof("  Query [%s], total fetched: %d", options.QueryText, len(commitItems))
	return commitItems, nil
}
