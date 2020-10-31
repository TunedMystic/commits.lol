package models

import (
	"html/template"
	"time"
)

// GitUser is the model for the git_user table.
type GitUser struct {
	ID        int    `db:"id"`
	Source    int    `db:"source"`
	Username  string `db:"username"`
	URL       string `db:"url"`
	AvatarURL string `db:"avatar_url"`
}

// GitRepo is the model for the git_repo table.
type GitRepo struct {
	ID          int    `db:"id"`
	Source      int    `db:"source"`
	Name        string `db:"name"`
	Description string `db:"description"`
	URL         string `db:"url"`
}

// GitCommit is the model for the git_commit table.
type GitCommit struct {
	ID              int       `db:"id"`
	Source          int       `db:"source"`
	AuthorID        int       `db:"author_id"`
	RepoID          int       `db:"repo_id"`
	Message         string    `db:"message"`
	MessageCensored string    `db:"message_censored"`
	SHA             string    `db:"sha"`
	URL             string    `db:"url"`
	Date            time.Time `db:"date"`
	CreatedAt       time.Time `db:"created_at"`
	Valid           bool      `db:"valid"`

	Author GitUser `db:"author"`
	Repo   GitRepo `db:"repo"`

	ColorBackground string
	ColorForeground string
}

// GitCommits is a slice of GitCommits.
type GitCommits []*GitCommit

// Term is the model for the term table.
type Term struct {
	ID   int    `db:"id"`
	Text string `db:"text"`
	Rank int    `db:"rank"`
}

// Terms is a slice of Term values.
type Terms []Term

// Strings converts the Terms slice into a slice of strings.
func (t Terms) Strings() []string {
	values := []string{}
	for _, term := range t {
		values = append(values, term.Text)
	}
	return values
}

// MessageCensoredHTML ...
func (c *GitCommit) MessageCensoredHTML() template.HTML {
	return template.HTML(c.MessageCensored)
}

// GetColorTheme ...
func (c *GitCommit) GetColorTheme() {
	// Colors ...
	colors := [][]string{

		// Ref: https://graf1x.com/shades-of-yellow-color-palette-chart/
		{"#fda50f", "#000000"}, // fire yellow
		{"#ffbf00", "#000000"}, // amber
		{"#fedc56", "#000000"}, // mustard
		{"#ffddaf", "#000000"}, // navajo
		{"#ffc30b", "#000000"}, // honey
		{"#ffd300", "#000000"}, // cyber
		{"#fada5e", "#000000"}, // royal
		{"#f8d373", "#000000"}, // mellow

		// Ref: https://www.eggradients.com/shades-of-green-color
		// {"#76ff7a", "#000000"}, // screamin green
		// {"#96ff36", "#000000"}, // spring green*
		{"#a7ff57", "#000000"}, // spring green*
		// {"#1fcecb", "#000000"}, // robin egg
		{"#0bda51", "#000000"}, // malachite
		// {"#24e860", "#000000"}, // lime green*

		// Ref: https://graf1x.com/shades-of-blue-color-palette/
		{"#73c2fb", "#000000"}, // maya
		{"#6593f5", "#000000"}, // cornflower
		// {"#074fbd", "#ffffff"}, // sapphire*
		{"#1f63ca", "#ffffff"}, // sapphire*
		{"#6cbff9", "#000000"}, // carolina*
		// {"#72cbf5", "#000000"}, // baby blue*

		// Ref: https://graf1x.com/24-shades-of-pink-color-palette/
		{"#fe7f9c", "#000000"}, // watermelon
		{"#ff66cc", "#000000"}, // rose pink
		{"#fb607f", "#000000"}, // brick

		// Ref: https://www.eggradients.com/shades-of-purple
		// {"#6147f1", "#ffffff"}, // electric indigo*
		{"#8a2be2", "#ffffff"}, // blue violet
		// {"#c71585", "#ffffff"}, // red violet
		// {"#6a5acd", "#ffffff"}, // slate blue
	}

	commitLength := len(c.Message) + len(c.Author.Username)
	colorIndex := commitLength % len(colors)
	color := colors[colorIndex]

	c.ColorBackground = color[0]
	c.ColorForeground = color[1]
}
