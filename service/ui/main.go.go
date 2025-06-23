package main

import (
    "fmt"
    "net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/html")
    fmt.Fprint(w, "<html><body><h1>OK</h1></body></html>")
}

func main() {
    http.HandleFunc("/", handler)
    http.ListenAndServe(":80", nil)
}