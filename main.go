package main

import (
	router "trade/Router"
	"trade/pkg/db"
)

func init() {

	db.Init()
}

func main() {
	router.Router()
}
