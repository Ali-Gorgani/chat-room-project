package handler

// ProtectedResponse represents the response structure for the protected endpoint.
type ProtectedResponse struct {
	Message string `json:"message" example:"You are authorized!"`
	Token   string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"`
}
