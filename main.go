package main

import "net/http"

func main() {
	mux := http.NewServeMux()

	http.ListenAndServe(":8080", mux)
}
