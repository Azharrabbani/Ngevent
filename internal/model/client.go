package model

type Client struct {
	IP        string `json:"ip"`
	UserAgent string `json:"user_agent"`
}
