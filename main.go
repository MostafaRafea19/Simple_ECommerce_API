package main

import (
	"api/database"
	"api/migrations"
	"api/routes"
	"log"
	"os"
)

func main() {
	service := os.Getenv("SERVICE")

	database.Connect()

	switch service {
	case "migration":
		if err := migrations.Migrate(database.GetDB()); err != nil {
			log.Fatalf("Could not run migrations: %v", err)
		}

		log.Println("Migrations ran successfully")

	case "customers", "sellers", "admins":
		r := routes.SetupRouter(service)
		port := os.Getenv("PORT")
		
		if err := r.Run(":" + port); err != nil {
			return
		}

	default:
		log.Fatal("Invalid service specified.")
		return
	}
}
