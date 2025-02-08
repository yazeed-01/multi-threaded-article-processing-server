package main

import (
	"mta/initializers"
	"mta/routes"
	"mta/services"
)

func init() {
	initializers.LoadEnvVar()
	initializers.ConnectDB()
}
func main() {

	initializers.InitElasticsearch()
	services.InitWorkerPool(5)
	r := routes.SetupRoutes()
	r.Run()
}
