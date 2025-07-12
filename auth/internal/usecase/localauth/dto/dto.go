package localauthdto

type LoginInput struct {
	Email    string             `json:"email"`
	Password string             `json:"password"`
	// Info     models.RequestInfo `json:"info"`
}

type SignupInput struct {
	Email    string             `json:"email"`
	Password string             `json:"password"`
	// Info     models.RequestInfo `json:"info"`
}
