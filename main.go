package main

import (
  "io"
  "log"
  "net/http"
  "github.com/gorilla/mux"
  //_ "github.com/lib/pq" //Ideally this will create a dependency on lib/pq
)

func hello(w http.ResponseWriter, req *http.Request) {
  name, ok := mux.Vars(req)["name"]
  if (!ok) {
    name = "world"
  }

  io.WriteString(w, "hello, "+name+"!\n")
}

func main() {
  // Define router
  router := mux.NewRouter()

  // Define routes here
  router.HandleFunc("/hello", hello)
  hello_router := router.PathPrefix("/hello").Subrouter()
  hello_router.HandleFunc("/", hello)
  hello_router.HandleFunc("/{name}", hello)

  // Mount router to path to be handled and start serving requests
  http.Handle("/", router)
  err := http.ListenAndServe(":12345", nil)
  if err != nil {
    log.Fatal("ListenAndServe: ", err)
  }
}
