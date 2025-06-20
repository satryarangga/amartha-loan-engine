package main

import (
	"flag"

	migration "github.com/satryarangga/amartha-loan-engine/database/migration"
)

func main() {
	flag.Parse()
	args := flag.Args()

	switch args[0] {
	case "migrate":
		migration.Migrate(args[1:])
	case "seed":
	}
}
