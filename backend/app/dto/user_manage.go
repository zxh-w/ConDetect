package dto

type UserMmanage struct {
	UID    string `json:"uid"`
	Name   string `json:"name"`
	Login  bool   `json:"login"`
	Docker bool   `json:"docker"`
}

type UserSearch struct {
	Name string `json:"name"`
}
