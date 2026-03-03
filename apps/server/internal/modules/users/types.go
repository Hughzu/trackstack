package users

type User struct {
	ID             string
	Email          string
	PasswordHash   string
	SessionVersion int
	CreatedAt      string
	LastLoginAt    *string
}
