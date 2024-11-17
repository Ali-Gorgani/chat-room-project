package user

type GetUserReq struct {
	Username string
}

type UserRes struct {
	ID       uint
	Username string
	Password string
	Email    string
	Role     Role
}

type Role struct {
	Name        string
	Premissions []string
}
