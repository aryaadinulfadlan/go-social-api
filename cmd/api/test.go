package main

import "net/http"

func (app *Application) Ping(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("PONG"))
}
func (app *Application) BasicAuthentication(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Authenticated as Basic Authentication"))
}
