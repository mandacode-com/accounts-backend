package logindto

type LocalLoginInput struct {
	Email    string             `json:"email"`
	Password string             `json:"password"`
	// Info     models.RequestInfo `json:"info"`
}
