package github

import (
	"testing"

	u "github.com/tunedmystic/commits.lol/app/utils"
)

func Test_Options_Empty(t *testing.T) {
	options := CommitSearchOptions{}
	u.AssertEqual(t, options.Empty(), true)

	options.QueryText = "lol"
	u.AssertEqual(t, options.Empty(), false)
}

func Test_Option_Serialize(t *testing.T) {
	type test struct {
		label    string
		opts     CommitSearchOptions
		expected string
	}

	tests := []test{
		{
			"blank",
			CommitSearchOptions{},
			"q=",
		},
		{
			"QueryText",
			CommitSearchOptions{
				QueryText: "lol",
			},
			"q='lol'",
		},
		{
			"QueryText_with_whitespace",
			CommitSearchOptions{
				QueryText: "awesome new feature",
			},
			"q='awesome+new+feature'",
		},
		{
			"date",
			CommitSearchOptions{
				FromDate: "2020-01-01",
				ToDate:   "2020-03-16",
			},
			"q=author-date:2020-01-01..2020-03-16",
		},
		{
			"date_FromDate_invalid",
			CommitSearchOptions{
				ToDate: "2020-03-16",
			},
			"q=",
		},
		{
			"date_ToDate_invalid",
			CommitSearchOptions{
				FromDate: "2020-01-01",
			},
			"q=",
		},
		{
			"Hash",
			CommitSearchOptions{
				Hash: "7f388fd42ab7d8342fbd0e0ece76a8505d228f1d",
			},
			"q=hash:7f388fd42ab7d8342fbd0e0ece76a8505d228f1d",
		},
		{
			"User",
			CommitSearchOptions{
				User: "TunedMystic",
			},
			"q=user:TunedMystic",
		},
		{
			"Org",
			CommitSearchOptions{
				Org: "golang",
			},
			"q=org:golang",
		},
		{
			"Repo",
			CommitSearchOptions{
				Repo: "tunedmystic/commits.lol",
			},
			"q=repo:tunedmystic/commits.lol",
		},
		{
			"Sort_asc",
			CommitSearchOptions{
				Sort: SortAsc,
			},
			"q=sort:author-date-asc",
		},
		{
			"Sort_desc",
			CommitSearchOptions{
				Sort: SortDesc,
			},
			"q=sort:author-date-desc",
		},
		{
			"example_query_1",
			CommitSearchOptions{
				QueryText: "should not be wrapped",
				FromDate:  "2020-09-27",
				ToDate:    "2020-10-01",
				User:      "rsc",
				Repo:      "golang/go",
				Sort:      SortDesc,
			},
			"q='should+not+be+wrapped'+author-date:2020-09-27..2020-10-01+user:rsc+repo:golang/go+sort:author-date-desc",
		},
		{
			"example_query_2",
			CommitSearchOptions{
				QueryText: "changelog++",
				FromDate:  "2020-01-01",
				ToDate:    "2020-03-16",
				Repo:      "hashicorp/vault",
				Page:      2,
			},
			"q='changelog%2B%2B'+author-date:2020-01-01..2020-03-16+repo:hashicorp/vault&page=2",
		},
	}

	for _, testItem := range tests {
		t.Run(testItem.label, func(t *testing.T) {
			u.AssertEqual(t, testItem.opts.Serialize(), testItem.expected)
		})
	}
}
