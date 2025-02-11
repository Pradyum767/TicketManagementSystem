package transport

import (
	service "go-jira/Service"
	"net/http"
)

func HttpHandler() {
	http.HandleFunc("/create", service.CreateIssue)
	http.HandleFunc("/updateStatus", service.UpdateStatus)
	http.HandleFunc("/updateAssignee", service.UpdateAssignee)
	http.HandleFunc("/addComment", service.AddComment)
}
