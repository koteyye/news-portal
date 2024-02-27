package resthandler

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/koteyye/news-portal/pkg/models"
)

func parseUserdata(r io.Reader) (*models.UserData, error) {
	var s models.UserData
	err := json.NewDecoder(r).Decode(&s)
	if err != nil {
		return nil, fmt.Errorf("decoding the userdata: %w", err)
	}
	if s.Login == "" {
		return nil, fmt.Errorf("login is empty")
	}
	if s.Password == "" {
		return nil, fmt.Errorf("password is empty")
	}
	return &s, nil
}