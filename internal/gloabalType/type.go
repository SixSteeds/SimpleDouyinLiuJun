package gloabalType

import "time"

type LoginSuccessMessage struct {
	IP        string    `json:"ip"`
	Logintime time.Time `json:"login_time"`
	UserId    int64     `json:"user_id"`
}

type UploadSuccessMessage struct {
	IP         string    `json:"ip"`
	Uploadtime time.Time `json:"upload_time"`
	UserId     int64     `json:"user_id"`
	PlayUrl    string    `json:"play_url"`
	DataLen    int64     `json:"data_len"`
}
