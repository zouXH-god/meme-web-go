package memes_cli

type MemeInfo struct {
	Key          string         `json:"key"`
	Params       MemeParams     `json:"params"`
	Keywords     []string       `json:"keywords"`
	Shortcuts    []MemeShortcut `json:"shortcuts"`
	Tags         []string       `json:"tags"`
	DateCreated  string         `json:"date_created"`
	DateModified string         `json:"date_modified"`
}

type MemeParams struct {
	MinImages    int          `json:"min_images"`
	MaxImages    int          `json:"max_images"`
	MinTexts     int          `json:"min_texts"`
	MaxTexts     int          `json:"max_texts"`
	DefaultTexts []string     `json:"default_texts"`
	Options      []MemeOption `json:"options"`
}

type MemeOption struct {
	Type        string      `json:"type"`
	Name        string      `json:"name"`
	Default     interface{} `json:"default"`
	Description *string     `json:"description"`
	ParserFlags ParserFlags `json:"parser_flags"`
	Choices     []string    `json:"choices"`
	Minimum     *float64    `json:"minimum"`
	Maximum     *float64    `json:"maximum"`
}

type ParserFlags struct {
	Short        bool     `json:"short"`
	Long         bool     `json:"long"`
	ShortAliases []string `json:"short_aliases"`
	LongAliases  []string `json:"long_aliases"`
}

type MemeShortcut struct {
	Pattern   string            `json:"pattern"`
	Humanized *string           `json:"humanized"`
	Names     []string          `json:"names"`
	Texts     []string          `json:"texts"`
	Options   map[string]string `json:"options"`
}

type Image struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}
