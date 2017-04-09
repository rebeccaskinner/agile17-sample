package main

import "github.com/julienschmidt/httprouter"

func main() {
	router := httprouter.New()
	router.Get("/old", handlers.dumpOld)
}
