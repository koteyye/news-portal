package models

import "github.com/gofrs/uuid"

type Profile struct {
	ID string `json:"userID,omitempty"`
	UserName string `json:"userName"`
	FirstName string `json:"firstName,omitempty"`
	LastName string `json:"lastName,omitempty"`
	SurName string `json:"surName,omitempty"`
	AvatarID uuid.UUID `json:"avatar,omitempty"`
	Roles []string `json:"roles,omitempty"`
}

type UserData struct {
	Login string `json:"login"`
	Password string `json:"password"`
	Profile *Profile `json:"profile,omitempty"`
}