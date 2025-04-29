package domain

type User struct {
	ID       uint
	Username string
	Email    string
	Password string
}

type UserRepository interface {
	Create(user *User) error
	GetByEmail(email string) (*User, error)
	UpdatePassword(email string, newPassword string) error
}

type AuthService interface {
	Signup(user *User) error
	Login(email, password string) (string, error)
	ConfirmPassword(email string) error
	ChangePassword(email, oldPwd, newPwd string) error
}
