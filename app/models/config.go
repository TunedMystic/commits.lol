package models

// BadWord is the model for the config_badword table.
type BadWord struct {
	ID   int    `db:"id"`
	Text string `db:"text"`
}

// GroupTerm is the model for the config_groupterm table.
type GroupTerm struct {
	ID    int    `db:"id"`
	Text  string `db:"text"`
	Group string `db:"groupname"`
}

// SearchTerm is the model for the config_searchterm table.
type SearchTerm struct {
	ID   int    `db:"id"`
	Text string `db:"text"`
	Rank int    `db:"rank"`
}

// BadWords is a slice of BadWord values.
type BadWords []BadWord

// ToStrings converts the BadWords slice into a slice of strings.
func (b BadWords) ToStrings() []string {
	values := make([]string, 0, len(b))
	for _, badword := range b {
		values = append(values, badword.Text)
	}
	return values
}

// GroupTerms is a slice of GroupTerm values.
type GroupTerms []GroupTerm

// ToMap converts the GroupTerms slice into a map of [string]string.
func (g GroupTerms) ToMap() map[string]string {
	values := make(map[string]string, len(g))
	for _, groupterm := range g {
		values[groupterm.Text] = groupterm.Group
	}
	return values
}

// SearchTerms is a slice of SearchTerm values.
type SearchTerms []SearchTerm

// ToStrings converts the SearchTerms slice into a slice of strings.
func (t SearchTerms) ToStrings() []string {
	values := make([]string, 0, len(t))
	for _, term := range t {
		values = append(values, term.Text)
	}
	return values
}
