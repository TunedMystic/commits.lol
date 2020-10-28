package github

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/tunedmystic/commits.lol/app/utils"
)

// These are qualifiers that are used to sort search results.
// Asc:  earliest date to latest date.
// Desc: latest date to earliest date.
const (
	SortAsc = iota + 1
	SortDesc
)

// CommitSearchOptions contains valid qualifiers / query params for the commit search endpoint.
type CommitSearchOptions struct {
	QueryText string
	FromDate  string // author date
	ToDate    string // author date
	Hash      string // commit hash
	User      string
	Org       string
	Repo      string // username/repo
	Sort      int
	Page      int
}

// IsEmpty checks if the options are empty.
func (opts CommitSearchOptions) IsEmpty() bool {
	return opts == (CommitSearchOptions{})
}

// Serialize ...
func (opts CommitSearchOptions) Serialize() string {
	queryString := fmt.Sprintf("q=%v", opts.queryParam())

	page := opts.pageParam()
	if page != "" {
		queryString += fmt.Sprintf("&page=%v", page)
	}

	return queryString
}

// PageParam ...
func (opts CommitSearchOptions) pageParam() string {
	if opts.Page == 0 {
		return ""
	}
	return strconv.Itoa(opts.Page)
}

// QueryParam ...
func (opts CommitSearchOptions) queryParam() string {
	qualifiers := []string{}

	if opts.QueryText != "" {
		qualifiers = append(qualifiers, "'"+url.QueryEscape(opts.QueryText)+"'")
	}

	if utils.IsValidDate(opts.FromDate) && utils.IsValidDate(opts.ToDate) {
		qualifiers = append(qualifiers, fmt.Sprintf("author-date:%v..%v", opts.FromDate, opts.ToDate))
	}

	if opts.Hash != "" {
		qualifiers = append(qualifiers, "hash:"+opts.Hash)
	}

	if opts.User != "" {
		qualifiers = append(qualifiers, "user:"+opts.User)
	}

	if opts.Org != "" {
		qualifiers = append(qualifiers, "org:"+opts.Org)
	}

	if opts.Repo != "" {
		qualifiers = append(qualifiers, "repo:"+opts.Repo)
	}

	if opts.Sort == SortAsc {
		qualifiers = append(qualifiers, "sort:author-date-asc")
	}

	if opts.Sort == SortDesc {
		qualifiers = append(qualifiers, "sort:author-date-desc")
	}

	return strings.Join(qualifiers, "+")
}
