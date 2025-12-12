package dtos

// Socket shared models
type Message struct {
	Type     string `json:"type"` // "join", "leave", "chat", "system"
	Room     string `json:"room"`
	Username string `json:"username"`
	Text     string `json:"text"`
	Time     string `json:"time"`
}
