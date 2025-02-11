package Service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	dataModel "myapp/DataModel"
	rep "myapp/Repository"
	"net/http"
	"net/url"

	"time"

	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/oauth2"
)

type Service struct {
	tpl               *template.Template
	googleOauthConfig *oauth2.Config
	Repository        rep.RepositoryInterface
}

var (
	oauthStateString = "pseudo-random"
	ticket           dataModel.Ticket
	picture          string
	resp             = make(map[string]string)
	token            *oauth2.Token
)

func NewService(tpl *template.Template, googleOauthConfig *oauth2.Config, repo rep.RepositoryInterface) ServiceInterface {

	return &Service{
		tpl:               tpl,
		googleOauthConfig: googleOauthConfig,
		Repository:        repo,
	}
}

func (s Service) GetTicket(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		cookie, err := r.Cookie("Logged-in")
		if err != nil || cookie.Value != token.AccessToken {
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}

		id := r.FormValue("id")

		ticket, err := s.Repository.GetTicket(id)
		if err != nil {

			fmt.Fprintln(w, err)
		} else {

			fmt.Println("Displayed Tickets: ", ticket)
			json.NewEncoder(w).Encode(ticket)

		}
	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("not found"))
	}
}

func (s Service) AllTicket(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		cookie, err := r.Cookie("Logged-in")
		if err != nil || cookie.Value != token.AccessToken {
			postBody, _ := json.Marshal(map[string]string{
				"message": "Login first ",
			})
			fmt.Fprintf(w, string(postBody))
			return
		}

		status := r.FormValue("status")

		cursor, err := s.Repository.Allticket(status)
		if err != nil {
			fmt.Fprintln(w, err)
		} else {
			var episodes []bson.M
			if err = cursor.All(context.TODO(), &episodes); err != nil {
				log.Fatal(err)
			}

			fmt.Println("Displayed all Tickets: ", episodes)
			json.NewEncoder(w).Encode(episodes)

		}
	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("not found"))
	}
}

func (s Service) Myticket(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		cookie, err := r.Cookie("Logged-in")
		if err != nil || cookie.Value != token.AccessToken {
			postBody, _ := json.Marshal(map[string]string{
				"message": "Login first ",
			})
			fmt.Fprintf(w, string(postBody))
			return
		}

		cursor, err := s.Repository.Myticket(ticket.Email)
		if err != nil {
			log.Fatal(err)
			fmt.Println(w, err)
		} else {
			var episodes []bson.M
			if err = cursor.All(context.TODO(), &episodes); err != nil {
				log.Fatal(err)
			}

			fmt.Println("Displayed all Tickets: ", episodes)
			json.NewEncoder(w).Encode(episodes)

		}
	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("not found"))
	}
}

func (s Service) AdminResponse(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		cookie, err := r.Cookie("Logged-in")
		if err != nil || cookie.Value != token.AccessToken {
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}
		id := r.FormValue("id")
		var ticket dataModel.Ticket
		ticket.Response = r.FormValue("response")
		ticket.Status = r.FormValue("status")
		ticket.Admincontact = r.FormValue("name")
		ticket2, result, err := s.Repository.AdminResponse(ticket, id)
		if err != nil {
			fmt.Println(w, "Error while Saving Reponse")
		} else {
			user, err := s.Repository.GetUser(ticket.Admincontact)
			form := url.Values{}
			form.Add("username", user.Email)
			form.Add("password", user.Password)
			form.Add("comment", ticket.Response)
			form.Add("ticketid", id)
			response, err := http.PostForm("http://localhost:8082/addComment", form)
			if err != nil {
				fmt.Println("The http request to update comment failed")
			}
			fmt.Println(" added comment   in jira", response, err)
			form.Add("status", ticket.Status)

			response, err = http.PostForm("http://localhost:8082/updateStatus", form)
			if err != nil {
				fmt.Println("The http request to update status failed")
			}

			fmt.Println("Change status in jira", response, err)
			form.Add("accountid", user.AccountID)

			response, err = http.PostForm("http://localhost:8082/updateAssignee", form)
			if err != nil {
				fmt.Println("The http request to update assignee failed")
			}
			fmt.Println("Changed assignee  in jira", response)
			fmt.Println("Updated a single document: ", ticket2, result)

			resp["email"] = ticket2.Email
			resp["user"] = "user"
			resp["message"] = "Response Recorded Succesfully"
			resp["picture"] = picture
			s.tpl.ExecuteTemplate(w, "Response.html", resp)
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("not found"))
	}
}

func (s Service) ReOpenTicket(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		cookie, err := r.Cookie("Logged-in")

		if err != nil || cookie.Value != token.AccessToken {
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}

		id := r.FormValue("id")
		ticket, result, err := s.Repository.ReOpenTicket(id)

		if err != nil {
			fmt.Fprintln(w, err)
		} else {

			form := url.Values{}
			user, err := s.Repository.GetUser(ticket.Email)

			form.Add("username", user.Email)
			form.Add("password", user.Password)
			form.Add("status", "Back to in progress")
			form.Add("ticketid", ticket.ID)
			form.Add("accountid", "")
			response, err := http.PostForm("http://localhost:8082/updateStatus", form)
			if err != nil {
				fmt.Println("The http request to update status failed")
			}
			form.Set("status", "Pending")
			response, err = http.PostForm("http://localhost:8082/updateStatus", form)
			if err != nil {
				fmt.Println("The http request to update status failed")
			}
			fmt.Println("Updates status to Pending", response, err)
			response, err = http.PostForm("http://localhost:8082/updateAssignee", form)
			if err != nil {
				fmt.Println("The http request to update assignee failed")
			}
			fmt.Println("Updated a single document: ", result)
			resp := make(map[string]string)
			resp["email"] = ticket.Email
			resp["user"] = "user"
			resp["message"] = "Ticket Reopened Successfully"
			resp["picture"] = picture
			s.tpl.ExecuteTemplate(w, "Response.html", resp)
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("not found"))
	}
}

func (s Service) CloseTicket(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		cookie, err := r.Cookie("Logged-in")
		if err != nil || cookie.Value != token.AccessToken {
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}
		id := r.FormValue("id")

		ticket, result, err := s.Repository.CloseTicket(id)

		if err != nil {

			log.Fatal(err)
			fmt.Fprintln(w, "Unable to Update Ticket Status")

		} else {
			form := url.Values{}
			user, _ := s.Repository.GetUser(ticket.Email)

			form.Add("username", user.Email)
			form.Add("password", user.Password)

			form.Add("ticketid", ticket.ID)
			if ticket.Status == "Resolve this issue" {

				form.Add("status", "Close")
				response, err := http.PostForm("http://localhost:8082/updateStatus", form)
				if err != nil {
					fmt.Println("The http request to update status failed")
				}
				fmt.Println("Successfull closed the ticket", response, err)
			} else {

				form.Add("status", "Cancel request")
				response, err := http.PostForm("http://localhost:8082/updateStatus", form)
				if err != nil {
					fmt.Println("The http request to update status failed")
				}
				form.Set("status", "Close")
				response, err = http.PostForm("http://localhost:8082/updateStatus", form)
				if err != nil {
					fmt.Println("The http request to update status failed")
				}
				fmt.Println("Successfull closed the ticket", response, err)
			}

			fmt.Println("Updated a single document: ", result)

			resp := make(map[string]string)

			resp["email"] = ticket.Email
			resp["user"] = "user"
			resp["message"] = "Ticket closed Succesfully"
			resp["picture"] = picture
			s.tpl.ExecuteTemplate(w, "Response.html", resp)
		}
	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("not found"))
	}
}
func (s Service) Raiserequest(w http.ResponseWriter, r *http.Request) {

	cookie, err := r.Cookie("Logged-in")
	if err != nil || cookie.Value != token.AccessToken {

		resp["message"] = "You need To Login First"
		s.tpl.ExecuteTemplate(w, "Response.html", resp)
		return
	}

	resp["email"] = ticket.Email
	resp["picture"] = picture
	s.tpl.ExecuteTemplate(w, "Raiserequest.html", resp)

}

func (s Service) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("Logged-in")
	if err != nil {
		return
	}
	cookie = &http.Cookie{

		Name:  "Logged-in",
		Value: "0",
	}

	http.SetCookie(w, cookie)
	for k := range resp {
		delete(resp, k)
	}

	postBody, _ := json.Marshal(map[string]string{
		"token": token.AccessToken,
	})
	responseBody := bytes.NewBuffer(postBody)
	response, err := http.Post("https://oauth2.googleapis.com/revoke", "application/json", responseBody)
	if err != nil {
		fmt.Println("Error while logging out")
	} else {
		defer response.Body.Close()
		fmt.Println("Logged Out Successfully")
		s.tpl.ExecuteTemplate(w, "index.html", nil)
	}

}
func (s Service) HandleMain(w http.ResponseWriter, r *http.Request) {

	s.tpl.ExecuteTemplate(w, "index.html", nil)

}

func (s Service) HandleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	fmt.Println("inside login")
	cookie, err := r.Cookie("Logged-in")
	if err != nil {
		cookie = &http.Cookie{
			Name:  "Logged-in",
			Value: "0",
		}
		http.SetCookie(w, cookie)
	}

	url := s.googleOauthConfig.AuthCodeURL(oauthStateString)

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)

}

func (s Service) HandleGoogleCallback(w http.ResponseWriter, r *http.Request) {

	var result map[string]interface{}
	fmt.Println("Inside callback")
	content, err := s.GetUserInfo(r.FormValue("state"), r.FormValue("code"))
	cookie, err := r.Cookie("Logged-in")
	if err == nil {
		cookie = &http.Cookie{

			Name:  "Logged-in",
			Value: token.AccessToken,
		}
		http.SetCookie(w, cookie)
	}

	if err != nil {

		fmt.Println(err.Error())

		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)

		return

	}

	json.Unmarshal([]byte(content), &result)

	picture = fmt.Sprintf("%v", result["picture"])
	if result["email"] == "pradyumsingh02@gmail.com" {

		fmt.Println("admin Logged in")
		ticket.Email = fmt.Sprintf("%v", result["email"])
		s.tpl.ExecuteTemplate(w, "AdminPage.html", result)

	} else {

		fmt.Println("User logged in")
		ticket.Email = fmt.Sprintf("%v", result["email"])

		s.tpl.ExecuteTemplate(w, "Login.html", result)
	}

}

func (s Service) GetUserInfo(state string, code string) ([]byte, error) {

	if state != oauthStateString {

		return nil, fmt.Errorf("invalid oauth state")

	}
	var err error
	token, err = s.googleOauthConfig.Exchange(context.Background(), code)

	if err != nil {

		return nil, fmt.Errorf("code exchange failed: %s", err.Error())

	}

	response, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {

		return nil, fmt.Errorf("failed getting user info: %s", err.Error())

	}

	defer response.Body.Close()

	contents, err := ioutil.ReadAll(response.Body)

	if err != nil {

		return nil, fmt.Errorf("failed reading response body: %s", err.Error())

	}

	return contents, nil

}

func (s Service) Queryas(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		cookie, err := r.Cookie("Logged-in")
		if err != nil || cookie.Value != token.AccessToken {
			resp["message"] = "You need to Login First"
			s.tpl.ExecuteTemplate(w, "Response.html", resp)
			return

		}

		dt := time.Now()
		ticket.Message = r.FormValue("message")
		ticket.Response = ""
		ticket.Status = "Pending"
		ticket.Time = dt.Format("01-02-2006 15:04:05")
		ticket.Admincontact = ""
		ticket.Description = r.FormValue("description")
		var rsp map[string]string
		form := url.Values{}
		user, err := s.Repository.GetUser(ticket.Email)

		form.Add("username", user.Email)
		form.Add("password", user.Password)
		form.Add("description", ticket.Description)
		form.Add("summary", ticket.Message)
		response, err := http.PostForm("http://localhost:8082/create", form)
		if err != nil {
			fmt.Println("The http request to create ticket in jira failed")
		}

		contents, err := ioutil.ReadAll(response.Body)
		json.Unmarshal([]byte(contents), &rsp)
		ticket.ID, _ = rsp["issueid"]
		form.Add("status", "Pending")
		form.Add("ticketid", ticket.ID)
		response, err = http.PostForm("http://localhost:8082/updateStatus", form)
		if err != nil {
			fmt.Println("The http request to update status failed")
		}
		fmt.Println("service", ticket, user, rsp)
		insertResult, err := s.Repository.Queryas(ticket)

		if err != nil {

			log.Fatal(err)

		}
		resp["message"] = "Your Ticket id is: " + ticket.ID
		resp["email"] = ticket.Email
		resp["user"] = "user"
		resp["picture"] = picture
		fmt.Println("Inserted a single document: ", insertResult)
		s.tpl.ExecuteTemplate(w, "Response.html", resp)
	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("not found"))
	}
}
