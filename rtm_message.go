package slackapi

// RTMMessage is a real-time message sent/received to/from Slack
type RTMMessage struct {
	ID      uint64 `json:"id"`
	Type    string `json:"type"`
	Channel string `json:"channel"`
	Text    string `json:"text"`
	User    string `json:"user"`
}

// RTMMessageResponse is received in response to RTMMessage requests
type RTMMessageResponse struct {
	OK        bool                `json:"ok"`       // success: true
	ReplyTo   int                 `json:"reply_to"` // ID of the request this responds to
	timestamp string              `json:"ts"`       // timestamp in the format "1355517523.000005"
	Text      string              `json:"text"`     // potentially modified version of the request
	Error     RTMMessageErrorCode `json:"error"`    // error, if OK==false
}

// RTMMessageErrorCode is included in RTMMessageResponse on error
type RTMMessageErrorCode struct {
	Code    int    `json:"code"`
	Message string `json:"msg"`
}
