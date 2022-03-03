package user

import (
	"errors"
	"log"
	"os"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dionaditya/go-production-ready-api/internal/models"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Service struct {
	DB *gorm.DB
}

type Token struct {
	Access_Token  string `json:",omitempty"`
	Refresh_Token string `json:",omitempty"`
}

//a struct to rep user account
type Payload struct {
	User models.User
	Token
}

type UserService interface {
	Register(user models.User) (Payload, error)
	Login(user string, password string) (Payload, error)
	UpdateUser(email string, newUsername string) (Payload, error)
	UpdatePassowrd(email, newPassword string) (bool, error)
	Validate(user models.User) (User error)
	RefreshToken(refreshToken string) (Payload, error)
}

func NewService(db *gorm.DB) *Service {
	return &Service{
		DB: db,
	}
}

func (s *Service) Validate(user models.User) (models.User, error) {

	if !strings.Contains(user.Email, "@") {
		return models.User{}, errors.New("Email address is required")
	}

	if len(user.Password) < 6 {
		return models.User{}, errors.New("Password less than 6 character")
	}

	var temp []models.User
	var result models.User
	//check for errors and duplicate emails
	err := s.DB.Find(&temp).Where("email = ?", user.Email).First(&result).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return models.User{}, errors.New("Connection error, please retry")
	}

	if result.Email != "" {
		return models.User{}, errors.New("Email already in use by another suer")
	}

	return user, nil
}

func (s *Service) Register(user models.User) (Payload, error) {

	if _, err := s.Validate(user); err != nil {
		return Payload{}, err
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)
	if result := s.DB.Save(&user); result.Error != nil {
		return Payload{}, result.Error
	}

	if user.ID <= 0 {
		return Payload{}, errors.New("Failed to create account, connection error")
	}

	user.Password = ""

	payload := Payload{User: user}
	return payload, nil
}

func (s *Service) Login(email, password string) (Payload, error) {
	var user []models.User
	var result models.User
	err := s.DB.Find(&user).Where("email = ?", email).First(&result).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return Payload{}, errors.New("Connection eror, please retry")
		}
		return Payload{}, errors.New("Email address not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(password))

	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		return Payload{}, errors.New("Invalid login credential")
	}

	result.Password = ""

	token, _ := GenerateJWT(result.ID, result.Email)
	payload := Payload{Token: token, User: result}
	return payload, nil
}

func (s *Service) RefreshToken(refreshToken string) (Payload, error) {
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Unexpected signing method")
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(os.Getenv("REFRESH_SECRET")), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Get the user record from database or
		// run through your business logic to verify if the user can log in

		var user []models.User
		var result models.User

		err := s.DB.Find(&user).Where("email = ?", claims["email"]).First(&result).Error

		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return Payload{}, errors.New("Connection eror, please retry")
			}
			return Payload{}, errors.New("Email address not found")
		}

		if claims["email"] == result.Email {

			newTokenPair, err := GenerateJWT(result.ID, result.Email)

			if err != nil {
				return Payload{}, err
			}

			return Payload{Token: newTokenPair, User: result}, nil
		}

		return Payload{}, err
	}

	return Payload{}, err
}

func GenerateJWT(userId uint, email string) (Token, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["client"] = userId
	claims["email"] = email
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()
	secretCode := []byte(os.Getenv("SECRET_CODE"))
	tokenString, err := token.SignedString(secretCode)

	if err != nil {
		log.Fatal("Failed to generate access token")
		return Token{}, err
	}

	refreshToken := jwt.New(jwt.SigningMethodHS256)
	rfClaims := refreshToken.Claims.(jwt.MapClaims)
	rfClaims["email"] = email
	rfClaims["exp"] = time.Now().Add(time.Minute * 24).Unix()

	rt, err := refreshToken.SignedString([]byte(os.Getenv("REFRESH_SECRET")))

	if err != nil {
		log.Fatal("Failed to generate refresh token")
		return Token{}, err
	}

	return Token{Access_Token: tokenString, Refresh_Token: rt}, nil
}

func (s *Service) GetUser(ID uint) (models.User, error) {
	var user models.User
	if result := s.DB.First(&user, ID); result.Error != nil {
		return models.User{}, result.Error
	}

	return user, nil
}

func (s *Service) UpdateUser(ID uint, newData struct{ Username string }) (models.User, error) {
	user, err := s.GetUser(ID)

	tempUser := user

	tempUser.Username = newData.Username

	if err != nil {
		return models.User{}, err
	}

	if result := s.DB.Model(&user).Updates(tempUser); result.Error != nil {
		return models.User{}, result.Error
	}

	return tempUser, nil
}
