package routers

import (
	"github.com/gorilla/mux"
	"net/http"
)

func Router() *mux.Router {
	r := mux.NewRouter()

	// static router: serve files under /static/<filename>
	r.PathPrefix("/static").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	//TODO: main entrypoint

	return r
}
