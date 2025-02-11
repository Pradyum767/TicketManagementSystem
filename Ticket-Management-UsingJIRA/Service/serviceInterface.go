package Service

import "net/http"

type ServiceInterface interface {
	GetTicket(w http.ResponseWriter, r *http.Request)
	Myticket(w http.ResponseWriter, r *http.Request)
	AdminResponse(w http.ResponseWriter, r *http.Request)
	ReOpenTicket(w http.ResponseWriter, r *http.Request)
	CloseTicket(w http.ResponseWriter, r *http.Request)
	Raiserequest(w http.ResponseWriter, r *http.Request)
	Logout(w http.ResponseWriter, r *http.Request)
	HandleMain(w http.ResponseWriter, r *http.Request)
	HandleGoogleLogin(w http.ResponseWriter, r *http.Request)
	HandleGoogleCallback(w http.ResponseWriter, r *http.Request)
	AllTicket(w http.ResponseWriter, r *http.Request)
	GetUserInfo(state string, code string) ([]byte, error)
	Queryas(w http.ResponseWriter, r *http.Request)
}
