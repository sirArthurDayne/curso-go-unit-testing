package main

import (
	"catching-pokemons/controller"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func Add(a, b int) int {
    return a + b
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/pokemon/{id}", controller.GetPokemon).Methods("GET")

	err := http.ListenAndServe(":5000", router)
	if err != nil {
		fmt.Print("Error found")
	}
}
