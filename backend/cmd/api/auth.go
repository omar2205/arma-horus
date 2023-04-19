package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"oskr.nl/arma-horus.go/internal/models"
	"oskr.nl/arma-horus.go/internal/utils"
)



func (app *application) SignWithGoogle(w http.ResponseWriter, r *http.Request) {
	config := &oauth2.Config{

		ClientID:     "137466250658-orfd65j7qtmup3bp0omfj20dpgm2oosj.apps.googleusercontent.com",
		ClientSecret: "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",
		RedirectURL:  "http://localhost:3000/v1/oauth/google",
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}

	state := uuid.NewV4().String()

	http.SetCookie(w, &http.Cookie{
		Name:  "oauthstate",
		Value: state,
		Path:  "/",
	})

	url := config.AuthCodeURL(state)

	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (app *application) HandleGoogleAuth(w http.ResponseWriter, r *http.Request) {
	config := &oauth2.Config{

		ClientID:     "137466250658-orfd65j7qtmup3bp0omfj20dpgm2oosj.apps.googleusercontent.com",
		ClientSecret: "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",
		RedirectURL:  "http://localhost:3000/v1/oauth/google",
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}

	coookie, err := r.Cookie("oauthstate")
	if err != nil || r.FormValue("state") != coookie.Value {
		app.writeJSON(w, 500, map[string]string{
			"error": "Invalid state!",
		}, nil)
		return
	}

	code := r.FormValue("code")
	token, err := config.Exchange(context.Background(), code)
	if err != nil {
		app.writeJSON(w, 400, map[string]string{
			"error": "Invalid code",
		}, nil)
		return
	}

	client := config.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		app.writeJSON(w, 500, map[string]string{
			"error": "Faield to get user info",
		}, nil)
		return
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		app.writeJSON(w, 500, map[string]string{
			"error": "Faield to get user info",
		}, nil)
		return
	}

	var userInfo map[string]interface{}
	err = json.Unmarshal(data, &userInfo)
	if err != nil {
		app.writeJSON(w, 500, map[string]string{
			"error": "Faield to parse user info",
		}, nil)
		return
	}
	email := userInfo["email"].(string)

	if !utils.Contains(app.config.Admins_email, email) {
		app.notPermittedResponse(w, r)
		return
	}

	user_token := "HORUS-" + uuid.NewV4().String()
	err = models.NewUser(email, user_token).Save()
	if err != nil {
		app.writeJSON(w, 500, map[string]string{
			"error": "Faield create user",
		}, nil)
		return
	}

	app.writeJSON(w, 201,
		map[string]string{"email": email, "token": user_token},
		nil,
	)
}
