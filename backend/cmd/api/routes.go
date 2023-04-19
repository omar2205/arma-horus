package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) HomePage(w http.ResponseWriter, r *http.Request,
	ps httprouter.Params) {

	w.Write([]byte(`
<h1>HORUS API</h1>
<a href="/v1/sign/google">Sign in with Google</a>
  `))
}

func (app *application) addServerHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Server", "ARMA/HORUS")
		next.ServeHTTP(w, r)
	})
}

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.GET("/", app.HomePage)
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodGet, "/v1/serverLogs", app.ServerLogsWS)

	router.HandlerFunc(http.MethodPost, "/v1/arma/start", app.ArmaServerStart)
	router.HandlerFunc(http.MethodPost, "/v1/arma/stop", app.ArmaServerStop)

	router.HandlerFunc(http.MethodGet, "/v1/sign/google", app.SignWithGoogle)
	router.HandlerFunc(http.MethodGet, "/v1/oauth/google", app.HandleGoogleAuth)

	router.HandlerFunc(http.MethodGet, "/v1/admin", app.authenticate(app.secretPage))

	return app.metrics(
		app.recoverPanic(
			app.enableCORS(app.addServerHeader(router)),
		),
	)
}

func (app *application) secretPage(w http.ResponseWriter, r *http.Request) {
	app.writeJSON(w, 200, map[string]string{"message": "super secret page"}, nil)
}
