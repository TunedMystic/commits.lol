package github

import (
	"encoding/json"
	"testing"
	"time"

	u "github.com/tunedmystic/commits.lol/app/utils"
)

func Test_CommitSearch_Unmarshal(t *testing.T) {
	response := CommitSearchResponse{}

	if err := json.Unmarshal([]byte(responseCommitSearch), &response); err != nil {
		t.Errorf("Failed to parse json: %v\n", err)
	}

	// Check CommitSearchResponse
	u.AssertEqual(t, response.TotalCount, 1)
	u.AssertEqual(t, len(response.CommitItems), 1)

	// Check CommitItem
	commitItem := response.CommitItems[0]
	u.AssertEqual(t, commitItem.URL, "https://github.com/TunedMystic/commits.lol/commit/7f388fd42ab7d8342fbd0e0ece76a8505d228f1d")
	u.AssertEqual(t, commitItem.SHA, "7f388fd42ab7d8342fbd0e0ece76a8505d228f1d")
	u.AssertEqual(t, commitItem.Score, 1.0)

	// Check Commit
	commit := commitItem.Commit
	u.AssertEqual(t, commit.Message, "Initial commit")
	u.AssertEqual(t, commit.Author.Date.Format(time.RFC3339), "2020-09-04T17:41:34-04:00")

	// Check Author
	author := commitItem.Author
	u.AssertEqual(t, author.Login, "TunedMystic")
	u.AssertEqual(t, author.AvatarURL, "https://avatars0.githubusercontent.com/u/6523726?v=4")
	u.AssertEqual(t, author.URL, "https://github.com/TunedMystic")

	// Check Repository
	repo := commitItem.Repo
	u.AssertEqual(t, repo.Name, "commits.lol")
	u.AssertEqual(t, repo.URL, "https://github.com/TunedMystic/commits.lol")
	u.AssertEqual(t, repo.Owner.Login, "TunedMystic")
}

func Test_CommitSearch_Unmarshal_no_author(t *testing.T) {
	response := CommitSearchResponse{}

	if err := json.Unmarshal([]byte(responseCommitSearchNoAuthor), &response); err != nil {
		t.Errorf("Failed to parse json: %v\n", err)
	}

	// Check Author
	u.AssertEqual(t, response.CommitItems[0].Author, Author{})
}

// ------------------------------------------------------------------
// Test Data
// ------------------------------------------------------------------

const responseCommitSearch = `{
    "total_count": 1,
    "incomplete_results": false,
    "items": [
        {
            "url": "https://api.github.com/repos/TunedMystic/commits.lol/commits/7f388fd42ab7d8342fbd0e0ece76a8505d228f1d",
            "sha": "7f388fd42ab7d8342fbd0e0ece76a8505d228f1d",
            "node_id": "MDY6Q29tbWl0MjkyOTQ2NDQ5OjdmMzg4ZmQ0MmFiN2Q4MzQyZmJkMGUwZWNlNzZhODUwNWQyMjhmMWQ=",
            "html_url": "https://github.com/TunedMystic/commits.lol/commit/7f388fd42ab7d8342fbd0e0ece76a8505d228f1d",
            "comments_url": "https://api.github.com/repos/TunedMystic/commits.lol/commits/7f388fd42ab7d8342fbd0e0ece76a8505d228f1d/comments",
            "commit": {
                "url": "https://api.github.com/repos/TunedMystic/commits.lol/git/commits/7f388fd42ab7d8342fbd0e0ece76a8505d228f1d",
                "author": {
                    "date": "2020-09-04T17:41:34.000-04:00",
                    "name": "Sandeep Jadoonanan",
                    "email": "someperson@gmail.com"
                },
                "committer": {
                    "date": "2020-09-04T17:41:34.000-04:00",
                    "name": "Sandeep Jadoonanan",
                    "email": "someperson@gmail.com"
                },
                "message": "Initial commit",
                "tree": {
                    "url": "https://api.github.com/repos/TunedMystic/commits.lol/git/trees/b2cf55f1573c8c3baf203acdc35b94d30c58ee76",
                    "sha": "b2cf55f1573c8c3baf203acdc35b94d30c58ee76"
                },
                "comment_count": 0
            },
            "author": {
                "login": "TunedMystic",
                "id": 6523726,
                "node_id": "MDQ6VXNlcjY1MjM3MjY=",
                "avatar_url": "https://avatars0.githubusercontent.com/u/6523726?v=4",
                "gravatar_id": "",
                "url": "https://api.github.com/users/TunedMystic",
                "html_url": "https://github.com/TunedMystic",
                "followers_url": "https://api.github.com/users/TunedMystic/followers",
                "following_url": "https://api.github.com/users/TunedMystic/following{/other_user}",
                "gists_url": "https://api.github.com/users/TunedMystic/gists{/gist_id}",
                "starred_url": "https://api.github.com/users/TunedMystic/starred{/owner}{/repo}",
                "subscriptions_url": "https://api.github.com/users/TunedMystic/subscriptions",
                "organizations_url": "https://api.github.com/users/TunedMystic/orgs",
                "repos_url": "https://api.github.com/users/TunedMystic/repos",
                "events_url": "https://api.github.com/users/TunedMystic/events{/privacy}",
                "received_events_url": "https://api.github.com/users/TunedMystic/received_events",
                "type": "User",
                "site_admin": false
            },
            "committer": {
                "login": "TunedMystic",
                "id": 6523726,
                "node_id": "MDQ6VXNlcjY1MjM3MjY=",
                "avatar_url": "https://avatars0.githubusercontent.com/u/6523726?v=4",
                "gravatar_id": "",
                "url": "https://api.github.com/users/TunedMystic",
                "html_url": "https://github.com/TunedMystic",
                "followers_url": "https://api.github.com/users/TunedMystic/followers",
                "following_url": "https://api.github.com/users/TunedMystic/following{/other_user}",
                "gists_url": "https://api.github.com/users/TunedMystic/gists{/gist_id}",
                "starred_url": "https://api.github.com/users/TunedMystic/starred{/owner}{/repo}",
                "subscriptions_url": "https://api.github.com/users/TunedMystic/subscriptions",
                "organizations_url": "https://api.github.com/users/TunedMystic/orgs",
                "repos_url": "https://api.github.com/users/TunedMystic/repos",
                "events_url": "https://api.github.com/users/TunedMystic/events{/privacy}",
                "received_events_url": "https://api.github.com/users/TunedMystic/received_events",
                "type": "User",
                "site_admin": false
            },
            "parents": [],
            "repository": {
                "id": 292946449,
                "node_id": "MDEwOlJlcG9zaXRvcnkyOTI5NDY0NDk=",
                "name": "commits.lol",
                "full_name": "TunedMystic/commits.lol",
                "private": false,
                "owner": {
                    "login": "TunedMystic",
                    "id": 6523726,
                    "node_id": "MDQ6VXNlcjY1MjM3MjY=",
                    "avatar_url": "https://avatars0.githubusercontent.com/u/6523726?v=4",
                    "gravatar_id": "",
                    "url": "https://api.github.com/users/TunedMystic",
                    "html_url": "https://github.com/TunedMystic",
                    "followers_url": "https://api.github.com/users/TunedMystic/followers",
                    "following_url": "https://api.github.com/users/TunedMystic/following{/other_user}",
                    "gists_url": "https://api.github.com/users/TunedMystic/gists{/gist_id}",
                    "starred_url": "https://api.github.com/users/TunedMystic/starred{/owner}{/repo}",
                    "subscriptions_url": "https://api.github.com/users/TunedMystic/subscriptions",
                    "organizations_url": "https://api.github.com/users/TunedMystic/orgs",
                    "repos_url": "https://api.github.com/users/TunedMystic/repos",
                    "events_url": "https://api.github.com/users/TunedMystic/events{/privacy}",
                    "received_events_url": "https://api.github.com/users/TunedMystic/received_events",
                    "type": "User",
                    "site_admin": false
                },
                "html_url": "https://github.com/TunedMystic/commits.lol",
                "description": "Spicy commits from across the web",
                "fork": false,
                "url": "https://api.github.com/repos/TunedMystic/commits.lol",
                "forks_url": "https://api.github.com/repos/TunedMystic/commits.lol/forks",
                "keys_url": "https://api.github.com/repos/TunedMystic/commits.lol/keys{/key_id}",
                "collaborators_url": "https://api.github.com/repos/TunedMystic/commits.lol/collaborators{/collaborator}",
                "teams_url": "https://api.github.com/repos/TunedMystic/commits.lol/teams",
                "hooks_url": "https://api.github.com/repos/TunedMystic/commits.lol/hooks",
                "issue_events_url": "https://api.github.com/repos/TunedMystic/commits.lol/issues/events{/number}",
                "events_url": "https://api.github.com/repos/TunedMystic/commits.lol/events",
                "assignees_url": "https://api.github.com/repos/TunedMystic/commits.lol/assignees{/user}",
                "branches_url": "https://api.github.com/repos/TunedMystic/commits.lol/branches{/branch}",
                "tags_url": "https://api.github.com/repos/TunedMystic/commits.lol/tags",
                "blobs_url": "https://api.github.com/repos/TunedMystic/commits.lol/git/blobs{/sha}",
                "git_tags_url": "https://api.github.com/repos/TunedMystic/commits.lol/git/tags{/sha}",
                "git_refs_url": "https://api.github.com/repos/TunedMystic/commits.lol/git/refs{/sha}",
                "trees_url": "https://api.github.com/repos/TunedMystic/commits.lol/git/trees{/sha}",
                "statuses_url": "https://api.github.com/repos/TunedMystic/commits.lol/statuses/{sha}",
                "languages_url": "https://api.github.com/repos/TunedMystic/commits.lol/languages",
                "stargazers_url": "https://api.github.com/repos/TunedMystic/commits.lol/stargazers",
                "contributors_url": "https://api.github.com/repos/TunedMystic/commits.lol/contributors",
                "subscribers_url": "https://api.github.com/repos/TunedMystic/commits.lol/subscribers",
                "subscription_url": "https://api.github.com/repos/TunedMystic/commits.lol/subscription",
                "commits_url": "https://api.github.com/repos/TunedMystic/commits.lol/commits{/sha}",
                "git_commits_url": "https://api.github.com/repos/TunedMystic/commits.lol/git/commits{/sha}",
                "comments_url": "https://api.github.com/repos/TunedMystic/commits.lol/comments{/number}",
                "issue_comment_url": "https://api.github.com/repos/TunedMystic/commits.lol/issues/comments{/number}",
                "contents_url": "https://api.github.com/repos/TunedMystic/commits.lol/contents/{+path}",
                "compare_url": "https://api.github.com/repos/TunedMystic/commits.lol/compare/{base}...{head}",
                "merges_url": "https://api.github.com/repos/TunedMystic/commits.lol/merges",
                "archive_url": "https://api.github.com/repos/TunedMystic/commits.lol/{archive_format}{/ref}",
                "downloads_url": "https://api.github.com/repos/TunedMystic/commits.lol/downloads",
                "issues_url": "https://api.github.com/repos/TunedMystic/commits.lol/issues{/number}",
                "pulls_url": "https://api.github.com/repos/TunedMystic/commits.lol/pulls{/number}",
                "milestones_url": "https://api.github.com/repos/TunedMystic/commits.lol/milestones{/number}",
                "notifications_url": "https://api.github.com/repos/TunedMystic/commits.lol/notifications{?since,all,participating}",
                "labels_url": "https://api.github.com/repos/TunedMystic/commits.lol/labels{/name}",
                "releases_url": "https://api.github.com/repos/TunedMystic/commits.lol/releases{/id}",
                "deployments_url": "https://api.github.com/repos/TunedMystic/commits.lol/deployments"
            },
            "score": 1.0
        }
    ]
}`

const responseCommitSearchNoAuthor = `{
    "total_count": 1,
    "incomplete_results": false,
    "items": [
        {
            "sha": "7f388fd42ab7d8342fbd0e0ece76a8505d228f1d",
            "html_url": "https://github.com/TunedMystic/commits.lol/commit/7f388fd42ab7d8342fbd0e0ece76a8505d228f1d",
            "commit": {
                "url": "https://api.github.com/repos/TunedMystic/commits.lol/git/commits/7f388fd42ab7d8342fbd0e0ece76a8505d228f1d",
                "author": null,
                "message": "Initial commit"
            },
            "author": null,
            "repository": {
                "name": "commits.lol",
                "owner": {
                    "login": "TunedMystic"
                },
                "html_url": "https://github.com/TunedMystic/commits.lol",
                "description": "Spicy commits from across the web"
            },
            "score": 1.0
        }
    ]
}`
