package main

import (
	"net/http"

	"nhooyr.io/websocket"
	"oskr.nl/arma-horus.go/internal/commander"
	"oskr.nl/arma-horus.go/internal/filewatcher"
)

func (app *application) ServerLogsWS(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, nil)
	if err != nil {
		app.logger.PrintError(err, map[string]string{
			"message": "error upgrading connection",
		})
	}
	defer conn.Close(websocket.StatusInternalError, "server error")

	watcher, err := filewatcher.NewFileWatcher(app.config.Server_folder + "server.logs")
	if err != nil {
		app.logger.PrintError(err, map[string]string{
			"message": "error watching server logs",
		})
	}

	content, err := watcher.ReadFile()
	if err != nil {
		app.logger.PrintError(err, map[string]string{
			"message": "error watching server logs",
		})
	}

	err = conn.Write(r.Context(), websocket.MessageText, content)
	if err != nil {
		app.logger.PrintError(err, map[string]string{
			"message": "error watching server logs",
		})
	}

	sub := watcher.Subscribe()
	defer watcher.Unsubscribe(sub)

	for {
		select {
		case <-r.Context().Done():
			conn.Close(websocket.StatusNormalClosure, "")
			return
		case line := <-sub:
			err = conn.Write(r.Context(), websocket.MessageText, []byte(line))
			if err != nil {
				app.logger.PrintError(err, map[string]string{
					"message": "error watching server logs",
				})
				return
			}
		}
	}
}

func (app *application) ArmaServerStart(w http.ResponseWriter, r *http.Request) {
	pid, err := commander.RunCMD(app.config.Server_script)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	app.writeJSON(w, 200, map[string]int{"pid": pid}, nil)
}

func (app *application) ArmaServerStop(w http.ResponseWriter, r *http.Request) {
	pid, err := commander.ReadPIDFromFile(app.config.Server_pid_file)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	err = commander.KillCMD(pid)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	app.writeJSON(w, 200, map[string]string{"message": "stopped server"}, nil)
}
