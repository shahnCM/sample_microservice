package migration

import (
	"auth_ms/pkg/model"
	"auth_ms/pkg/provider/database/mariadb10"
	"log"
	"os"
)

func RunMigration() {
	runMigration := os.Getenv("RUN_MIGRATION")

	if runMigration == "TRUE" {
		db := mariadb10.GetMariaDb10()

		if err := db.AutoMigrate(&model.User{}, &model.Session{}, &model.Token{}); err != nil {
			log.Fatal("Failed to migrate schema:", err)
		}
		log.Println("Database migration completed")
	} else {
		log.Println("RUN_MIGRATION is not set to true, skipping migration")
	}

}
