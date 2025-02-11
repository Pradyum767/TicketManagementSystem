package dataModel

type User struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	AccountID string `json:"accountid"`
	Address   string `json:"address"`
	Contact   int    `json:"contact"`
}
