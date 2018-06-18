package main

import (
	"flag"
	"log"
	"os"

	"github.com/pressly/goose"

	"github.com/khisakuni/schubert-services/user/database"

	// Init DB drivers.
	_ "github.com/lib/pq"

	_ "github.com/khisakuni/schubert-services/user/migrations"
)

var (
	flags = flag.NewFlagSet("goose", flag.ExitOnError)
	dir   = flags.String("dir", ".", "directory with migration files")
)

func main() {
	flags.Usage = usage
	flags.Parse(os.Args[1:])

	args := flags.Args()

	if len(args) > 1 && args[0] == "create" {
		if err := goose.Run("create", nil, *dir, args[1:]...); err != nil {
			log.Fatalf("goose run: %v", err)
		}
		return
	}

	if len(args) < 1 {
		flags.Usage()
		return
	}

	command := args[0]
	if command == "-h" || command == "--help" {
		flags.Usage()
		return
	}

	if err := goose.SetDialect("postgres"); err != nil {
		log.Fatal(err)
	}

	db, err := database.New()
	checkError(err)

	arguments := []string{}
	if len(args) > 3 {
		arguments = append(arguments, args[1:]...)
	}

	if err := goose.Run(command, db, *dir, arguments...); err != nil {
		log.Fatalf("goose run: %v", err)
	}
}

func usage() {
	log.Print(usagePrefix)
	flags.PrintDefaults()
	log.Print(usageCommands)
}

var (
	usagePrefix = `Usage: goose [OPTIONS] DRIVER DBSTRING COMMAND
Drivers:
    postgres
Examples:
    goose postgres "user=postgres dbname=postgres sslmode=disable" status
Options:
`

	usageCommands = `
Commands:
    up                   Migrate the DB to the most recent version available
    up-to VERSION        Migrate the DB to a specific VERSION
    down                 Roll back the version by 1
    down-to VERSION      Roll back to a specific VERSION
    redo                 Re-run the latest migration
    status               Dump the migration status for the current DB
    version              Print the current version of the database
    create NAME [sql|go] Creates new migration file with next version
`
)

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
