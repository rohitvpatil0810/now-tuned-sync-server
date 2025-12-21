package model

type ClientType string

const (
	ClientChrome  ClientType = "chrome-extension"
	ClientAndroid ClientType = "android"
	ClientIOS     ClientType = "ios"
)

type ClientInfo struct {
	ClientID   string     `json:"clientId"` // stable UUID per install
	ClientType ClientType `json:"clientType"`
	DeviceID   string     `json:"deviceId"`        // optional (android/ios)
	TabID      *int       `json:"tabId,omitempty"` // only for chrome
}
