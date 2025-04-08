package main

import "net/http"

func (app *Application) Test(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}
