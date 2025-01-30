package main

import (
	"github.com/okyws/go-banking/routes"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	routes.StartServer()
}
