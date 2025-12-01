package im

type LoginRequest struct {
	UserID string `json:"user_id"`
}

type LoginResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type SendRequest struct {
	From    string `json:"from"`    // sender
	To      string `json:"to"`      // receiver
	Content string `json:"content"` // message content
}

// SendResponse
type SendResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// Message
type Message struct {
	From      string `json:"from"`
	To        string `json:"to"`
	Content   string `json:"content"`
	Timestamp int64  `json:"timestamp"`
}

type QueryUserRequest struct {
	UserID string `json:"user_id"`
}

type QueryUserResponse struct {
	UserID   string `json:"user_id"`
	Location string `json:"location"`
	Online   bool   `json:"online"`
}

// ClientMessage is the message sent by the client to the server.
type ClientMessage struct {
	Type    string `json:"type"`
	To      string `json:"to"`
	Content string `json:"content"`
}

// ServerMessage is the message sent by the server to the client.
type ServerMessage struct {
	Type    string `json:"type"`
	From    string `json:"from"`
	Content string `json:"content"`
}
