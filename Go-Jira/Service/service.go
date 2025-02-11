package service

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/andygrunwald/go-jira"
)

func CreateIssue(w http.ResponseWriter, r *http.Request) {

	jiraClient := getJiraClient(r.FormValue("username"), r.FormValue("password"))
	i := jira.Issue{
		Fields: &jira.IssueFields{

			Description: r.FormValue("description"),
			Type: jira.IssueType{
				Name: "[System] Service request",
			},
			Project: jira.Project{
				Key: "TMP",
			},
			Summary: r.FormValue("summary"),
		},
	}

	issue, _, err := jiraClient.Issue.Create(&i)
	if err != nil {
		panic(err)
	}
	fmt.Println("Issue created->\n", issue)
	fmt.Println(issue.ID)
	postBody, _ := json.Marshal(map[string]string{
		"issueid": issue.ID,
	})
	fmt.Println(postBody)
	w.Write(postBody)

}

func UpdateStatus(w http.ResponseWriter, r *http.Request) {
	jiraClient := getJiraClient(r.FormValue("username"), r.FormValue("password"))
	var transitionID string
	possibleTransitions, _, _ := jiraClient.Issue.GetTransitions(r.FormValue("ticketid"))
	for _, v := range possibleTransitions {
		fmt.Println(v.Name)
		if v.Name == r.FormValue("status") {
			transitionID = v.ID
			break
		}
	}

	resp, err := jiraClient.Issue.DoTransition(r.FormValue("ticketid"), transitionID)
	fmt.Println(resp, err)
}

func AddComment(w http.ResponseWriter, r *http.Request) {

	jiraClient := getJiraClient(r.FormValue("username"), r.FormValue("password"))
	issue, _, _ := jiraClient.Issue.Get(r.FormValue("ticketid"), nil)
	jiraClient.Issue.AddComment(issue.ID, &jira.Comment{Body: r.FormValue("comment")})

}

func UpdateAssignee(w http.ResponseWriter, r *http.Request) {

	jiraClient := getJiraClient(r.FormValue("username"), r.FormValue("password"))
	issue, _, _ := jiraClient.Issue.Get(r.FormValue("ticketid"), nil)
	fmt.Printf("%s: %+v\n", issue.Key, issue.Fields.Summary)

	assignee := &jira.User{AccountID: r.FormValue("accountid")}
	resp, err := jiraClient.Issue.UpdateAssignee(issue.ID, assignee)
	fmt.Println(resp, err)
}

func getJiraClient(username string, password string) *jira.Client {
	tp := jira.BasicAuthTransport{

		Username: username,
		Password: password,
	}

	jiraClient, _ := jira.NewClient(tp.Client(), "https://golangteamb.atlassian.net/")
	return jiraClient

}
