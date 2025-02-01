package main

import (
	"github.com/rismapa/go-banking/routes"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	routes.StartServer()
}
