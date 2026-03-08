package auth

import (
	"fmt"
	"strings"

	"github.com/DiegoAmin/AmazonClone_PAP/internal/logger"
)

// User represents a user in the system with a username, password, and role (e.g., "admin" or "customer").
type User struct {
	Username string
	Password string
	Role     string
}

// AuthStore manages user authentication and stores user information in memory.
type AuthStore struct {
	Users map[string]*User // key: username, value: User completo
}

// NewAuthStore initializes and returns a new AuthStore instance with default users.
func NewAuthStore() *AuthStore {
	a := &AuthStore{
		Users: make(map[string]*User),
	}

	// Default users so the system can be tested without registering
	a.Users["admin"] = &User{Username: "admin", Password: "admin123", Role: "admin"}
	a.Users["carlos"] = &User{Username: "carlos", Password: "carlos123", Role: "customer"}
	a.Users["diego"] = &User{Username: "diego", Password: "diego123", Role: "customer"}
	a.Users["diego2"] = &User{Username: "diego2", Password: "diego2123", Role: "customer"}

	return a
}

// Register adds a new user to the AuthStore. It returns an error if the username already exists or if the role is invalid.
func (a *AuthStore) Register(username, password, role string) error {
	username = strings.ToLower(username)
	if _, exists := a.Users[username]; exists {
		logger.Log(fmt.Sprintf("ERROR: user %s already exists", username))
		return fmt.Errorf("user %s already exists", username)
	}

	if role != "admin" && role != "customer" {
		logger.Log(fmt.Sprintf("ERROR: invalid role: %s", role))
		return fmt.Errorf("invalid role: %s", role)
	}

	a.Users[username] = &User{
		Username: username,
		Password: password,
		Role:     role,
	}
	logger.Log(fmt.Sprintf("AUTH: user registered: %s, role: %s", username, role))
	return nil
}

// Login validates the username and password and returns the user if successful.
func (a *AuthStore) Login(username, password string) (*User, error) {
	username = strings.ToLower(username)
	if user, exists := a.Users[username]; !exists {
		logger.Log(fmt.Sprintf("ERROR: user %s not found", username))
		return nil, fmt.Errorf("user %s not found", username)
	} else if user.Password != password {
		logger.Log(fmt.Sprintf("ERROR: invalid password for user %s", username))
		return nil, fmt.Errorf("invalid password for user %s", username)
	} else {
		logger.Log(fmt.Sprintf("AUTH: user logged in: %s, role: %s", username, user.Role))
		return user, nil
	}
}

func (a *AuthStore) ListUsers() []*User {
	logger.Log("ADMIN: user list requested")
	users := make([]*User, 0, len(a.Users))
	for _, user := range a.Users {
		users = append(users, user)
	}
	return users
}
