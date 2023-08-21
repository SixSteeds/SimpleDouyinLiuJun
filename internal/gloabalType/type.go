package gloabalType

import "time"

type LoginSuccessMessage struct {
	IP        string    `json:"ip"`
	Logintime time.Time `json:"login_time"`
	UserId    int64     `json:"user_id"`
}
