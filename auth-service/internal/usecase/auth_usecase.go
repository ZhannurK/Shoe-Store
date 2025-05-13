package usecase

import (
	"auth-service/internal/entities"
	"auth-service/internal/repositories"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"log"
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
	GetUserByID(id string) (*entities.User, error)
}

type Cache interface {
	Get(key string) (string, error)
	Set(key string, value string, ttl time.Duration) error
	Del(key string) error
}

type AuthUsecase struct {
	repo   repositories.UserRepository
	jwtKey []byte
	cache  Cache
}

func NewAuthUseCase(repo *repositories.UserMongoRepo, jwtKey []byte, cache Cache) *AuthUsecase {
	return &AuthUsecase{repo: repo, jwtKey: jwtKey, cache: cache}
}

func (uc *AuthUsecase) GetUserByID(id string) (*entities.User, error) {
	cacheKey := "user:" + id
	if cached, _ := uc.cache.Get(cacheKey); cached != "" {
		var user entities.User
		_ = json.Unmarshal([]byte(cached), &user)
		return &user, nil
	}

	user, err := uc.repo.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	data, _ := json.Marshal(user)
	err = uc.cache.Set(cacheKey, string(data), 5*time.Minute)
	if err != nil {
		return nil, err
	}

	if cached, _ := uc.cache.Get(cacheKey); cached != "" {
		log.Println("✅ [CACHE HIT] Returning cached value for", cacheKey)
		return user, nil
	} else {
		log.Println("❌ [CACHE MISS] Fetching from DB and caching", cacheKey)
	}

	return user, nil
}

func (uc *AuthUsecase) SignUp(ctx context.Context, email, password, name string) (string, error) {
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

//func (uc *AuthUsecase) ConfirmEmail(ctx context.Context, token string) error {
//	user, err := uc.repo.FindByToken(ctx, token)
//	if err != nil {
//		return err
//	}
//	user.Verified = true
//	user.ConfirmationToken = ""
//	return uc.repo.Update(ctx, user)
//}

func (uc *AuthUsecase) ConfirmEmail(ctx context.Context, token string) error {
	user, err := uc.repo.FindByToken(ctx, token)
	if err != nil {
		return err
	}
	user.Verified = true
	user.ConfirmationToken = ""

	err = uc.repo.Update(ctx, user)
	if err != nil {
		return err
	}

	err = uc.cache.Del("user:" + user.ID.Hex())
	if err != nil {
		return err
	}

	data, _ := json.Marshal(user)
	err = uc.cache.Set("user:"+user.Email, string(data), 0)
	if err != nil {
		return err
	}
	return nil
}

func (uc *AuthUsecase) Login(ctx context.Context, email, password string) (string, *entities.User, error) {
	cacheKey := "user:" + email
	var user *entities.User
	if cached, _ := uc.cache.Get(cacheKey); cached != "" {
		log.Println("✅ [CACHE HIT] Returning cached value for", cacheKey)
		_ = json.Unmarshal([]byte(cached), &user)
	} else {
		log.Println("❌ [CACHE MISS] Fetching from DB and caching", cacheKey)
		var err error
		user, err = uc.repo.FindByEmail(ctx, email)
		if err != nil {
			return "", nil, errors.New("invalid credentials")
		}
		data, _ := json.Marshal(user)
		err = uc.cache.Set(cacheKey, string(data), 5*time.Minute)
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

//func (uc *AuthUsecase) ChangePassword(ctx context.Context, email, oldPwd, newPwd string) error {
//	user, err := uc.repo.FindByEmail(ctx, email)
//	if err != nil {
//		return err
//	}
//	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPwd)); err != nil {
//		return errors.New("incorrect old password")
//	}
//	newHashed, err := hashPassword(newPwd)
//	if err != nil {
//		return err
//	}
//	user.Password = newHashed
//	return uc.repo.Update(ctx, user)
//}

func (uc *AuthUsecase) ChangePassword(ctx context.Context, email, oldPwd, newPwd string) error {
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

	err = uc.repo.Update(ctx, user)
	if err != nil {
		return err
	}

	err = uc.cache.Del("user:" + user.Email)
	if err != nil {
		return err
	}

	data, _ := json.Marshal(user)
	err = uc.cache.Set("user:"+user.Email, string(data), 0)
	if err != nil {
		return err
	}
	return nil
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
