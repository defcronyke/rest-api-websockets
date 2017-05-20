package RestApiWebsockets

import (
    "net/http"
    "github.com/gorilla/mux"
    "github.com/rs/cors"
)

func init() {
    r := mux.NewRouter()
    apiPath := "/api"
    auth := &Auth{}
    r.HandleFunc(apiPath + "/createAccount", auth.CreateAccount).Methods("POST")
    r.HandleFunc(apiPath + "/login", auth.Login).Methods("POST")
    r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static")))
    http.Handle("/", cors.Default().Handler(r))
}
