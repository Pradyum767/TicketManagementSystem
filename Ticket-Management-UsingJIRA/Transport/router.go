package Transport

import (
	svc "myapp/Service"
	"net/http"
)

func HttpRouter(service svc.ServiceInterface) {

	http.HandleFunc("/", service.HandleMain)
	http.HandleFunc("/login", service.HandleGoogleLogin)
	http.Handle("/stuff/", http.StripPrefix("/stuff", http.FileServer(http.Dir("./WebPages/assets/"))))
	http.HandleFunc("/callback", service.HandleGoogleCallback)
	http.HandleFunc("/send", service.Queryas)
	http.HandleFunc("/myticket", service.Myticket)
	http.HandleFunc("/getticket", service.GetTicket)
	http.HandleFunc("/allticket", service.AllTicket)
	http.HandleFunc("/Raiserequest", service.Raiserequest)
	// http.HandleFunc("/adminTickets", service.AdminTickets)
	http.HandleFunc("/adminResponse", service.AdminResponse)
	http.HandleFunc("/close", service.CloseTicket)
	http.HandleFunc("/reopen", service.ReOpenTicket)
	http.HandleFunc("/logout", service.Logout)

}
