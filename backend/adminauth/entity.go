package adminauth

type LoginRequest struct {
	Password string
}

type LoginResponse struct {
	Token string
}
