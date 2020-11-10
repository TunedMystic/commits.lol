package models

import (
	"html/template"
	"time"

	"github.com/tunedmystic/commits.lol/app/utils"
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
	Group           string    `db:"groupname"`
	ColorBackground string    `db:"color_bg"`
	ColorForeground string    `db:"color_fg"`

	Author GitUser `db:"author"`
	Repo   GitRepo `db:"repo"`
}

// GitCommits is a slice of GitCommits.
type GitCommits []GitCommit

// MessageCensoredHTML ...
func (c *GitCommit) MessageCensoredHTML() template.HTML {
	return template.HTML(c.MessageCensored)
}

// SetCensoredMessage cleans the commit message and sets it as the `MessageCensored` field.
// Returns true if message was censored.
// Returns false if there were no bad words to be cleaned.
func (c *GitCommit) SetCensoredMessage(cl utils.Cleaner) bool {
	cleanedMsg, wordsCensored := cl.Clean(c.Message)

	// If the cleaned message is the same as the commit's message, then nothing was cleaned.
	// If there were no words censored, then nothing was cleaned.
	// In any of these cases, return false to express that the Commit was not updated.
	if cleanedMsg == c.Message || wordsCensored == 0 {
		return false
	}

	c.MessageCensored = cleanedMsg
	return true
}

// SetGroup assigns the commit to a group based on the commit message.
func (c *GitCommit) SetGroup(g utils.Grouper) bool {
	commitGroup := g.Group(c.Message)

	// If the generated group is the same as the commit's group, then nothing was changed.
	// In that case, return false to express that the Commit was not updated.
	if commitGroup == c.Group {
		return false
	}

	c.Group = commitGroup
	return true
}

// SetColorTheme sets the background and foreground color based on
// various attributes of the given Commit.
func (c *GitCommit) SetColorTheme() {
	// Colors ...
	colors := [][]string{

		// Ref: https://graf1x.com/shades-of-yellow-color-palette-chart/
		// {"#fda50f", "#000000"}, // fire yellow
		{"#ffbf00", "#000000"}, // amber
		{"#fedc56", "#000000"}, // mustard
		{"#ffddaf", "#000000"}, // navajo
		// {"#ffc30b", "#000000"}, // honey
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
		// {"#6593f5", "#000000"}, // cornflower
		// {"#074fbd", "#ffffff"}, // sapphire*
		// {"#1f63ca", "#ffffff"}, // sapphire*
		// {"#6cbff9", "#000000"}, // carolina*
		// {"#72cbf5", "#000000"}, // baby blue*

		// Ref: https://graf1x.com/24-shades-of-pink-color-palette/
		// {"#fe7f9c", "#000000"}, // watermelon
		// {"#ff66cc", "#000000"}, // rose pink
		// {"#fb607f", "#000000"}, // brick
		{"#cea8ff", "#000000"}, // light purple*
		{"#ff9ff3", "#000000"}, // jigglypuff

		// Ref: https://www.eggradients.com/shades-of-purple
		// {"#6147f1", "#ffffff"}, // electric indigo*
		// {"#8a2be2", "#ffffff"}, // blue violet
		// {"#c71585", "#ffffff"}, // red violet
		// {"#6a5acd", "#ffffff"}, // slate blue
		// {"#d4449a", "#ffffff"}, // slate pink*
		{"#ce3876", "#ffffff"}, // slate pink*

		// canada
		// {"	#1dd1a1", "#000000"}, // wild caribbean green
		// {"#ffc93c", "#000000"}, // wild caribbean green
		// {"#ff6b6b", "#ffffff"}, // pastel red
		// {"#222f3e", "#ffffff"}, // imperial primer
		// {"#feca57", "#000000"}, // casandora yellow
		// {"#ff9ff3", "#000000"}, // jigglypuff
		// {"#ff9f43", "#000000"}, // double dragon skin
		// {"#0ca9f2", "#000000"}, // jade dust
		// {"#9a89f7", "#000000"}, // joust blue

		// // spanish
		// {"#ff793f", "#000000"}, // synthetic pumpkin
		// {"#2b5cb7", "#ffffff"}, // c64 purple
		// {"#2c2c54", "#ffffff"}, // lucky point

		// // india
		// {"#2c3a47", "#ffffff"}, // ship's officer
		// {"#b33771", "#ffffff"}, // fiery fuchsia
		// {"#fd7272", "#000000"}, // georgia peach
		// {"#1B9CFC", "#ffffff"}, // clear chill

	}

	commitLength := len(c.Message) + len(c.Author.Username)
	colorIndex := commitLength % len(colors)
	color := colors[colorIndex]

	c.ColorBackground = color[0]
	c.ColorForeground = color[1]
}
