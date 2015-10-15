package slackapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// ChannelInfo holds a text value set by a user - used in Channel
type ChannelInfo struct {
	Value   string `json:"value"`
	Creator string `json:"creator"`
	LastSet uint64 `json:"last_set"`
}

// Channel represents a Slack channel
type Channel struct {
	ID                 string      `json:"id"`
	Name               string      `json:"name"`
	IsChannel          bool        `json:"is_channel"`
	Created            uint64      `json:"created"`
	Creator            string      `json:"creator"`
	IsArchived         bool        `json:"is_archived"`
	IsGeneral          bool        `json:"is_general"`
	Members            []string    `json:"members"`
	IsMember           bool        `json:"is_member"`
	LastRead           string      `json:"last_read"`
	UnreadCount        uint64      `json:"unread_count"`
	UnreadCountDisplay uint64      `json:"unread_count_display"`
	Topic              ChannelInfo `json:"topic"`
	Purpose            ChannelInfo `json:"purpose"`
	// TODO: `json:"latest"`
}

// IsDirectMessageChannel returns whether the channel ID is formatted as a direct message channel
func IsDirectMessageChannel(channelID string) bool {
	return strings.HasPrefix(channelID, "D")
}

// GetChannelInfo returns a *Channel from the input channelID
func GetChannelInfo(token string, channelID string) (*Channel, error) {
	if IsDirectMessageChannel(channelID) {
		return nil, fmt.Errorf("channelID '%s' is not a channel, but a direct message")
	}

	url := fmt.Sprintf("https://slack.com/api/channels.info?token=%s&channel=%s", token, channelID)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Error requesting info for channel %s: %s", channelID, err)
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API request failed with code %d", resp.StatusCode)
	}

	// get the response from HTTP body
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("Error received trying to close the channel info body: %s", err)
	}

	// unmarshal response
	var channelInfo Channel
	err = json.Unmarshal(body, &channelInfo)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshalling channel info response: %s", err)
	}

	return &channelInfo, nil
}
