package main

import "github.com/voyagerstudio/haiku-auth/pkg/api"

func main() {
	srv := api.NewServer("", 8080)
	srv.ListenAndServe()
}
