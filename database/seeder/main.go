package seeder

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"github.com/satryarangga/amartha-loan-engine/config"
)

func Seed() {

	db, err := config.InitDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	seederDir := "database/seeder/sql"

	// Read all SQL files in the seeder directory
	files, err := ioutil.ReadDir(seederDir)
	if err != nil {
		log.Fatal("Failed to read seeder directory:", err)
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".sql") {
			filePath := filepath.Join(seederDir, file.Name())
			log.Printf("Executing SQL file: %s", file.Name())

			// Read the SQL file content
			content, err := ioutil.ReadFile(filePath)
			if err != nil {
				log.Printf("Failed to read file %s: %v", file.Name(), err)
				continue
			}

			// Split the content by semicolon to get individual SQL statements
			statements := strings.Split(string(content), ";")

			for i, statement := range statements {
				statement = strings.TrimSpace(statement)
				if statement == "" {
					continue
				}

				// Execute each SQL statement
				result := db.Exec(statement)
				if result.Error != nil {
					log.Printf("Failed to execute statement %d in %s: %v", i+1, file.Name(), result.Error)
					continue
				}

				log.Printf("Successfully executed statement %d in %s", i+1, file.Name())
			}
		}
	}

	log.Println("Seeder executed successfully")
}
