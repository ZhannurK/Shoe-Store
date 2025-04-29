package usecase

import (
	"auth-service/internal/entities"
	"auth-service/internal/repositories"
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type AuthUseCase interface {
	SignUp(ctx context.Context, email, password, name string) (string, error)
	ConfirmEmail(ctx context.Context, token string) error
	Login(ctx context.Context, email, password string) (string, *entities.User, error)
	ChangePassword(ctx context.Context, email, oldPwd, newPwd string) error
}

type authUseCase struct {
	repo   repositories.UserRepository
	jwtKey []byte
}

func NewAuthUseCase(repo repositories.UserRepository, jwtKey []byte) AuthUseCase {
	return &authUseCase{repo: repo, jwtKey: jwtKey}
}

func (uc *authUseCase) SignUp(ctx context.Context, email, password, name string) (string, error) {
	if _, err := uc.repo.FindByEmail(ctx, email); err == nil {
		return "", errors.New("user already exists")
	}

	hashedPassword, err := hashPassword(password)
	if err != nil {
		return "", err
	}

	token, err := generateRandomToken()
	if err != nil {
		return "", err
	}

	user := entities.User{
		ID:                primitive.NewObjectID(),
		Email:             email,
		Name:              name,
		Password:          hashedPassword,
		Verified:          false,
		ConfirmationToken: token,
		Role:              "user",
	}

	if err := uc.repo.Create(ctx, &user); err != nil {
		return "", err
	}

	return token, nil
}

func (uc *authUseCase) ConfirmEmail(ctx context.Context, token string) error {
	user, err := uc.repo.FindByToken(ctx, token)
	if err != nil {
		return err
	}
	user.Verified = true
	user.ConfirmationToken = ""
	return uc.repo.Update(ctx, user)
}

func (uc *authUseCase) Login(ctx context.Context, email, password string) (string, *entities.User, error) {
	user, err := uc.repo.FindByEmail(ctx, email)
	if err != nil {
		return "", nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", nil, errors.New("invalid credentials")
	}

	if !user.Verified {
		return "", nil, errors.New("email not confirmed")
	}

	claims := jwt.MapClaims{
		"email": user.Email,
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString(uc.jwtKey)
	if err != nil {
		return "", nil, err
	}
	return tokenStr, user, nil
}

func (uc *authUseCase) ChangePassword(ctx context.Context, email, oldPwd, newPwd string) error {
	user, err := uc.repo.FindByEmail(ctx, email)
	if err != nil {
		return err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPwd)); err != nil {
		return errors.New("incorrect old password")
	}
	newHashed, err := hashPassword(newPwd)
	if err != nil {
		return err
	}
	user.Password = newHashed
	return uc.repo.Update(ctx, user)
}

func generateRandomToken() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}
