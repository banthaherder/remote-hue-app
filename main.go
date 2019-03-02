package main

import "net/http"

func main() {
	http.HandleFunc("/", handler)
	if err := http.ListenAndServe(":80", nil); err != nil {
		panic(err)
	}
}
