package auth

type SignUpRequest struct {
	Email    string
	Username string
	Password string
}

type SignUpResponse struct {
	Token string
}

type LoginRequest struct {
	Username string
	Password string
}

type LoginResponse struct {
	Token string
}
