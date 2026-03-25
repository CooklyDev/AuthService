package presentation

// SessionIDResponse documents successful register and login responses.
type SessionIDResponse struct {
	Success bool   `json:"success" example:"true"`
	Data    string `json:"data" format:"uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
}

// LogoutResponse documents a successful logout response.
type LogoutResponse struct {
	Success bool `json:"success" example:"true"`
}

// ErrorResponse documents failed responses that contain error details and no data.
type ErrorResponse struct {
	Success bool      `json:"success" example:"false"`
	Error   ErrorInfo `json:"error"`
}

// ResolveSessionData documents the payload returned by session resolution.
type ResolveSessionData struct {
	SessionID string `json:"SessionID" format:"uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	UserID    string `json:"UserID" format:"uuid" example:"7b1f5cde-8a0f-4bfb-b5e8-9d4dd0c5a7f3"`
}

// ResolveSessionResponse documents a successful resolve response.
type ResolveSessionResponse struct {
	Success bool               `json:"success" example:"true"`
	Data    ResolveSessionData `json:"data"`
}
