package models

type Friend struct {
	UserInfo
	Message string `json:"message"`
	MsgType int64  `json:"msg_type"`
}
