package models

type User struct {
	UserId int "json:userId"
	UserPwd string `json:"userPwd"`
	UserName string `json:"userName"`
}

