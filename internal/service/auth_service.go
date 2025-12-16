package service

import (
	"errors"
	"time"

	"github.com/JuanTobonV/blog_app/internal/model"
	"github.com/JuanTobonV/blog_app/internal/store"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)


type IAuthService interface {
	Register(username, password string) (*model.User, error)
	Login(username, login string) (string, error) //returns JWT token
	ValidateToken(tokenString string) (*jwt.Token, error)
	GetUserIDFromToken(tokenString string) (int, error)

}

type authService struct {
	userStore store.IUserStore  // Dependency on user store
    jwtSecret []byte            // Secret key for JWT
}

func NewAuthService(userStore store.IUserStore, jwtSecret string) IAuthService {
	return &authService{
		userStore: userStore,
		jwtSecret: []byte(jwtSecret),
	}
}

func (s *authService) Register(username, password string) (*model.User, error) {
	// 1. Validate input

	if username == "" || password == "" {
		return nil, errors.New("Username and password are required")
	}

	// 2. check if user already exists
	existingUser, err := s.userStore.GetByUsername(username)

	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("User already exists")
	}

	// 3. Hash the password using bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return nil, err
	}

	// 4. Create user object with hashed password
	user := &model.User{
		Username: username,
		Password: string(hashedPassword),
	}


    // 5. Save to database using store

	createdUser, err := s.userStore.Create(user)

	createdUser.Password = ""

	return createdUser, nil

}

func (s *authService) Login(username, password string) (string, error) {
	// 1. Find user by username

	user, err := s.userStore.GetByUsername(username)

	if err != nil {
		return "", err
	}

	if user == nil {
		return "", errors.New("Invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil {
		return "", errors.New("invalid credentials")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":user.Id,
		"username": user.Username,
		"exp": time.Now().Add(time.Hour * 24).Unix(), // Expires in 24 hours
	})

	tokenString, err := token.SignedString(s.jwtSecret)

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *authService) ValidateToken(tokenString string) (*jwt.Token, error) {
    // Parse and validate the token
    token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
        // Verify signing method
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, errors.New("invalid signing method")
        }
        return s.jwtSecret, nil
    })
    
    if err != nil {
        return nil, err
    }
    
    // Check if token is valid
    if !token.Valid {
        return nil, errors.New("invalid token")
    }
    
    return token, nil
}

func (s *authService) GetUserIDFromToken(tokenString string) (int, error) {
    token, err := s.ValidateToken(tokenString)
    if err != nil {
        return 0, err
    }
    
    // Extract claims
    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok {
        return 0, errors.New("invalid token claims")
    }
    
    // Get user_id from claims
    userID, ok := claims["user_id"].(float64) // JSON numbers are float64
    if !ok {
        return 0, errors.New("user_id not found in token")
    }
    
    return int(userID), nil
}
