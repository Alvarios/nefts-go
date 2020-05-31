package config

type Error struct {
	Code    int    `json:"status"`
	Message string `json:"message"`
}
