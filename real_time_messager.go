package slackapi

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/websocket"
	"io/ioutil"
	"net/http"
	"sync/atomic"
)

// rtmStartResponse is received from Slack on rtm.start request
type rtmStartResponse struct {
	OK    bool   `json:"ok"`
	URL   string `json:"url"`
	Error string `json:"error"`
}

// RealTimeMessager sends/receives real-time messages from Slack
type RealTimeMessager struct {
	token               string
	lastRequestID       uint64 // the last request ID sent to Slack
	webSocketConnection *websocket.Conn
}

// NewRealTimeMessager returns a new RealTimeMessager
func NewRealTimeMessager(token string) (*RealTimeMessager, error) {
	return &RealTimeMessager{
		token: token,
	}, nil
}

// newRequestID returns the next available request id
func (rtm *RealTimeMessager) newRequestID() uint64 {
	return atomic.AddUint64(&rtm.lastRequestID, 1)
}

// Connect connects to the Slack real-time Messaging endpoint
func (rtm *RealTimeMessager) Connect() error {
	// make a request to the rtm.start endpoint
	url := fmt.Sprintf("https://slack.com/api/rtm.start?token=%s", rtm.token)
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("Error received making a HTTP GET request to Slack's rtm.start endpoint: %s", err)
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("API request failed with code %d", resp.StatusCode)
	}

	// get response from HTTP body
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return fmt.Errorf("Error received trying to close the rtm.start body: %s", err)
	}

	// unmarshal response
	var rtmResponse rtmStartResponse
	err = json.Unmarshal(body, &rtmResponse)
	if err != nil {
		return fmt.Errorf("Error unmarshalling rtm.start response: %s", err)
	}

	// Check if OK
	if !rtmResponse.OK {
		return fmt.Errorf("Slack rtm.start response indicated *NOT* okay,  error: %s", rtmResponse.Error)
	}

	// dial the web socket
	wsConn, err := websocket.Dial(rtmResponse.URL, "", "https://api.slack.com/")
	if err != nil {
		return fmt.Errorf("Failure dialing web socket at %s: %s", rtmResponse.URL, err)
	}
	rtm.webSocketConnection = wsConn
	return nil
}

// ReceiveMessage blocks while waiting for the next real-time message from Slack
func (rtm *RealTimeMessager) ReceiveMessage() (*RTMMessage, error) {
	var message RTMMessage

	if err := websocket.JSON.Receive(rtm.webSocketConnection, &message); err != nil {
		return nil, fmt.Errorf("Error received when trying to fetch a JSON message from Slack RTM web socket connection: %s", err)
	}
	return &message, nil
}

// SendMessage sends a message request to Slack
func (rtm *RealTimeMessager) SendMessage(request RTMMessage) error {
	request.ID = rtm.newRequestID()
	return websocket.JSON.Send(rtm.webSocketConnection, request)
}
