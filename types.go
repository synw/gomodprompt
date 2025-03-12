package modprompt

type SpacingSlots struct {
	System    *int `json:"system,omitempty"`
	User      *int `json:"user,omitempty"`
	Assistant *int `json:"assistant,omitempty"`
}

type PromptBlock struct {
	Schema  string  `json:"schema"`
	Message *string `json:"message,omitempty"`
}

type TurnBlock struct {
	User      string `json:"user"`
	Assistant string `json:"assistant"`
	Tool      string `json:"tool,omitempty"`
}

type ToolDefinition struct {
	Def      string `json:"def,omitempty"`
	Call     string `json:"call,omitempty"`
	Response string `json:"response,omitempty"`
}

type LmTemplate struct {
	ID         string          `json:"id"`
	Name       string          `json:"name"`
	User       string          `json:"user"`
	Assistant  string          `json:"assistant"`
	System     *PromptBlock    `json:"system,omitempty"`
	Shots      []TurnBlock     `json:"shots,omitempty"`
	Stop       []string        `json:"stop,omitempty"`
	Linebreaks *SpacingSlots   `json:"linebreaks,omitempty"`
	AfterShot  *string         `json:"afterShot,omitempty"`
	Prefix     *string         `json:"prefix,omitempty"`
	ToolsDef   *ToolDefinition `json:"tools,omitempty"` // New field
}

type PromptTemplate struct {
	ID              string
	Name            string
	User            string
	Assistant       string
	System          *PromptBlock
	Shots           []TurnBlock
	Stop            []string
	Linebreaks      *SpacingSlots
	AfterShot       *string
	Prefix          *string
	_extraSystem    string
	_replaceSystem  string
	_extraAssistant string
	_replacePrompt  string
	history         []HistoryTurn
	Tools           []map[string]interface{} `json:"tools,omitempty"`     // New field
	ToolsDef        *ToolDefinition          `json:"tools_def,omitempty"` // New field
}

type ImgData struct {
	ID   int    `json:"id"`
	Data string `json:"data"`
}

type HistoryTurn struct {
	User      string    `json:"user"`
	Assistant string    `json:"assistant"`
	Tool      string    `json:"tool,omitempty"`
	Images    []ImgData `json:"images,omitempty"`
}
