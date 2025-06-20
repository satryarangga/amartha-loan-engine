package migrations

import (
	"log"

	_ "github.com/lib/pq"
	"github.com/pressly/goose"
	"github.com/satryarangga/amartha-loan-engine/internal/app"
)

/*
currently hardcoded for mysql
*/
func Migrate(args []string) {
	if len(args) < 1 {
		log.Fatalf("missing argument: ./{bin-file} [goose-command]")
		return
	}
	config, err := app.NewConfig()
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	//reading all custom-defined args
	migrationDir, gooseCommand := "database/migration/sql", args[0]

	//check db connection
	db, err := goose.OpenDBWithDriver(config.DBDriver, config.DBConnection)
	if err != nil {
		log.Fatalf("goose: failed to open DB: %v\n", err)
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Fatalf("goose: failed to close DB: %v\n", err)
		}
	}()

	//include additional goose pre-defined args
	arguments := []string{}
	if len(args) > 3 {
		arguments = append(arguments, args[3:]...)
	}

	//executing goose actual
	if err1 := goose.Run(gooseCommand, db, migrationDir, arguments...); err1 != nil {
		log.Fatalf("goose %v: %v", gooseCommand, err1)
	}
}
