package usecase

import (
	"Shoe-Store/internal/user/domain"
	"Shoe-Store/pkg/token"

	"golang.org/x/crypto/bcrypt"
)

type authService struct {
	repo domain.UserRepository
}

func NewAuthService(r domain.UserRepository) domain.AuthService {
	return &authService{repo: r}
}

func (s *authService) Signup(user *domain.User) error {
	hash, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	user.Password = string(hash)
	return s.repo.Create(user)
}

func (s *authService) Login(email, password string) (string, error) {
	user, err := s.repo.GetByEmail(email)
	if err != nil {
		return "", err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", err
	}
	return token.GenerateToken(user.ID)
}

func (s *authService) ConfirmPassword(email string) error {
	_, err := s.repo.GetByEmail(email)
	return err
}

func (s *authService) ChangePassword(email, oldPwd, newPwd string) error {
	user, err := s.repo.GetByEmail(email)
	if err != nil {
		return err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPwd)); err != nil {
		return err
	}
	hash, _ := bcrypt.GenerateFromPassword([]byte(newPwd), 14)
	return s.repo.UpdatePassword(email, string(hash))
}
