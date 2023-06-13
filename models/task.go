package models

type Task struct {
	Id     int    `json:"id"`
	UserId int    `json:"-"`
	Title  string `json:"title"`
	Body   string `json:"desc"`
	Status string `json:"status"`
}
