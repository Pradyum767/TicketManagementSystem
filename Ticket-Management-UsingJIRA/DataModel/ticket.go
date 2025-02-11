package dataModel

type Ticket struct {
	ID           string `json:"id,omitempty"`
	Email        string `json:"email"`
	Message      string `json:"message"`
	Description  string `json:"description"`
	Response     string `json:"response"`
	Status       string `json:"status"`
	Time         string `json:"time"`
	Admincontact string `json:"admincontact"`
}
