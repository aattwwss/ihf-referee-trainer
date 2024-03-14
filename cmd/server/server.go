package main

import (
	"log/slog"
	"net/http"
)

func main() {
	//start http server
	http.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.Handle("GET /test", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	}))

	err := http.ListenAndServe(":1234", nil)
	if err != nil {
		slog.Error("http server error", slog.String("error", err.Error()))
		panic(err)
	}
}
