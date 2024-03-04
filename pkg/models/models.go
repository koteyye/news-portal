package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"

	"github.com/gofrs/uuid"
)

const DefaultRole = "reader" // DefaultRole роль назначаемая по умолчанию новому пользователю.

// Profile профиль пользователя
type Profile struct {
	ID        string    `json:"userID,omitempty"`
	UserName  string    `json:"userName"`
	FirstName string    `json:"firstName,omitempty"`
	LastName  string    `json:"lastName,omitempty"`
	SurName   string    `json:"surName,omitempty"`
	AvatarID  uuid.UUID `json:"avatar,omitempty"`
	Roles     []string  `json:"roles,omitempty"`
}

// UserData пользовательские данные
type UserData struct {
	Login    string   `json:"login"`
	Password string   `json:"password"`
	Profile  *Profile `json:"profile,omitempty"`
}

// ParseUserData сериализует UserData
func ParseUserData(r io.Reader) (*UserData, error) {
	var s UserData
	err := json.NewDecoder(r).Decode(&s)
	if err != nil {
		return nil, fmt.Errorf("decoding the userdata: %w", err)
	}
	if s.Login == "" {
		return nil, errors.New("login is empty")
	}
	if s.Password == "" {
		return nil, errors.New("password is empty")
	}
	return &s, nil
}

// ParseProfile сериализует Profile
func ParseProfile(r io.Reader) (*Profile, error) {
	var p Profile
	err := json.NewDecoder(r).Decode(&p)
	if err != nil {
		return nil, fmt.Errorf("decoding the profile: %w", err)
	}
	if p.ID == "" {
		return nil, errors.New("userID is empty")
	}
	if p.UserName == "" {
		return nil, errors.New("username is empty")
	}
	return &p, nil
}