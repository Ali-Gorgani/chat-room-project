package domain

type User struct {
	ID       int
	Username string
	Password string
	Email    string
	Role     Role
	Profile  Profile
}

type Role struct {
	ID          int
	Name        string
	Premissions []string
}

type Profile struct {
	ID             int
	FirstName      string
	LastName       string
	ProfilePicture string
}

type Auth struct {
	AccessToken string
}
