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

type NewsAttributes struct {
	ID          string   `json:"news_id,omitempty"`
	Title       string   `json:"title"`
	Author      *Profile `json:"author,omitempty"`
	Description string   `json:"description"`
	Content     *File    `json:"content,omitempty"`
	Preview     *File    `json:"preview,omitempty"`
	State       string   `json:"state,omitempty"`
	CreatedAt   string   `json:"createdAt,omitempty"`
	UpdatedAt   string   `json:"updatedAt,omitempty"`
	UserCreated *Profile `json:"user_created,omitempty"`
	UserUpdated *Profile `json:"user_updated,omitempty"`
}

type File struct {
	ID         string `json:"file_id"`
	MimeType   string `json:"mime_type"`
	BucketName string `json:"bucket_name"`
	FileName   string `json:"file_name"`
}

type Like struct {
	ID        string   `json:"id"`
	Liker     *Profile `json:"liker"`
	CreatedAt string   `json:"created_at"`
	UpdatedAt string   `json:"update_at"`
}

type Comment struct {
	ID          string   `json:"id"`
	Author      *Profile `json:"author"`
	CreatedAt   string   `json:"createdAt"`
	UpdatedAt   string   `json:"updatedAt"`
	TextComment string   `json:"textComment"`
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

func ParseComment(r io.Reader) (*Comment, error) {
	var c Comment
	err := json.NewDecoder(r).Decode(&c)
	if err != nil {
		return nil, fmt.Errorf("decoding the comment: %w", err)
	}
	if c.TextComment == "" {
		return nil, errors.New("text comment is empty")
	}
	return &c, nil
}
