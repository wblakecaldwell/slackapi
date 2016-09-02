package slackapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type userResponse struct {
	Ok    bool   `json:"ok"`
	Error string `json:"error"`
	User  User   `json:"user"`
}

// User represents a Slack user
type User struct {
	ID      string `json:"id"`
	TeamID  string `json:"team_id"`
	Name    string `json:"name"`
	Deleted bool   `json:"deleted"`
	Color   string `json:"color"`
	Profile struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		RealName  string `json:"real_name"`
		Email     string `json:"email"`
		Skype     string `json:"skype"`
		Phone     string `json:"phone"`
		Image24   string `json:"image_24"`
		Image32   string `json:"image_32"`
		Image48   string `json:"image_48"`
		Image72   string `json:"image_72"`
		Image192  string `json:"image_192"`
	} `json:"profile"`
	IsAdmin  bool `json:"is_admin"`
	IsOwner  bool `json:"is_owner"`
	Has2FA   bool `json:"has_2va"`
	HasFiles bool `json:"has_files"`
}

// GetUserInfo returns a *User from the input userID
func GetUserInfo(token string, userID string) (*User, error) {
	url := fmt.Sprintf("https://slack.com/api/users.info?token=%s&user=%s", token, userID)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Error requesting info for user %s: %s", userID, err)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API request failed with code %d", resp.StatusCode)
	}

	// get the response from HTTP body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error received trying to close the user info body: %s", err)
	}
	defer resp.Body.Close()

	// unmarshal response
	var userResponse userResponse
	err = json.Unmarshal(body, &userResponse)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshalling user info response: %s", err)
	}

	if !userResponse.Ok {
		return nil, fmt.Errorf("Response from Slack was not okay - error: %s", userResponse.Error)
	}

	return &userResponse.User, nil
}
