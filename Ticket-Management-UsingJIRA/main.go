package main

import (
	"context"
	"fmt"
	"html/template"

	rep "myapp/Repository"
	svc "myapp/Service"
	Transport "myapp/Transport"
	"net/http"

	"golang.org/x/oauth2"

	"golang.org/x/oauth2/google"
)

var (
	googleOauthConfig *oauth2.Config
	tpl               *template.Template
)

func init() {
	tpl = template.Must(template.ParseGlob("WebPages/templates/*.html"))
	googleOauthConfig = &oauth2.Config{

		RedirectURL:  "http://localhost:8080/callback",
		ClientID:     "30246144593-khhu80p0undgb5911o5le0lv7jkekcet.apps.googleusercontent.com",
		ClientSecret: "wyXXzd_eLm6nbDaZDLYvyP3c",
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}

}

func main() {

	fmt.Println(googleOauthConfig.ClientID)
	repo := rep.NewRepository()
	var service svc.ServiceInterface = svc.NewService(tpl, googleOauthConfig, &repo)
	defer repo.Client.Disconnect(context.TODO())
	Transport.HttpRouter(service)
	http.ListenAndServe(":8080", nil)

}
