package auth

import (
	"encoding/json"
	"fmt"
	"os"
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

	// Save the updated users to the JSON file after every registration.
	if err := a.Save("users.json"); err != nil {
		logger.Log(fmt.Sprintf("ERROR: failed to save users after registration: %s", err.Error()))
	}
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

// ListUsers returns a slice of all users in the AuthStore. This function is intended for admin use only.
func (a *AuthStore) ListUsers() []*User {
	logger.Log("ADMIN: user list requested")
	users := make([]*User, 0, len(a.Users))
	for _, user := range a.Users {
		users = append(users, user)
	}
	return users
}

// Save saves the AuthStore data to a JSON file. It returns an error if the file cannot be written.
func (a *AuthStore) Save(filename string) error {
	data, err := json.MarshalIndent(a, "", "  ")
	if err != nil {
		logger.Log(fmt.Sprintf("ERROR: failed to marshal users: %s", err.Error()))
		return err
	}
	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		logger.Log(fmt.Sprintf("ERROR: failed to save users to file: %s", err.Error()))
		return err
	}
	logger.Log(fmt.Sprintf("AUTH: users saved to file: %s", filename))
	return nil
}

// Load reads the AuthStore data from a JSON file and returns a new AuthStore instance with the loaded data.
func Load(filename string) (*AuthStore, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		logger.Log(fmt.Sprintf("ERROR: failed to read users from file: %s", err.Error()))
		return nil, err
	}
	var authStore AuthStore
	err = json.Unmarshal(data, &authStore)
	if err != nil {
		logger.Log(fmt.Sprintf("ERROR: failed to unmarshal users from file: %s", err.Error()))
		return nil, err
	}
	logger.Log(fmt.Sprintf("AUTH: users loaded from file: %s", filename))
	return &authStore, nil
}
