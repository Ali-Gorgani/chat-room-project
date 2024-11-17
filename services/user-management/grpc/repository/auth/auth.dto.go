package auth

type HashPasswordReq struct {
	Password string
}

type HashPasswordRes struct {
	HashedPassword string
}

type VerifyTokenReq struct {
	Token string
}

type VerifyTokenRes struct {
	ID       int
	Username string
	Email    string
	Role     string
}
