package models

type AdminDashLog struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}
