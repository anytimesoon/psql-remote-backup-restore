package main

import (
	"database/sql"
	"fmt"
	pg "github.com/habx/pg-commands"
	"log"
	"net/url"
)

func restore(restoreFile string) {
	localDb := &pg.Postgres{
		Host:     local.host,
		Port:     local.port,
		DB:       local.name,
		Username: local.user,
		Password: local.password,
	}

	log.Println("Preparing for dropping of local schema")
	connString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", localDb.Username, url.QueryEscape(localDb.Password), localDb.Host, localDb.Port, localDb.DB)
	db, err := sql.Open("postgres", connString)
	if err != nil {
		log.Println("Failed to prepare for drop")
		log.Fatal(err)
	}

	log.Println("Dropping local schema")
	err = spinner.Start()
	spinner.Message("executing")
	query := "DROP SCHEMA public CASCADE; CREATE SCHEMA public; GRANT ALL ON SCHEMA public TO public;"
	_, err = db.Exec(query)
	if err != nil {
		log.Println("Failed to drop")
		log.Fatal(err)
	}
	err = spinner.Stop()
	if err != nil {
		log.Println("Failed to stop spinner. Continuing with restore...", err)
	}
	log.Println("Drop successful")

	log.Println("Preparing restore")
	pg.PGRestoreCmd = pathToRestore
	restore, err := pg.NewRestore(localDb)
	if err != nil {
		log.Println("Failed preparing restore")
		log.Fatal(err)
	}

	restore.Role = local.user
	restore.Options = restoreOptions

	log.Println("Starting restore")
	err = spinner.Start()
	spinner.Message("restoring data")
	restoreExec := restore.Exec(restoreFile, pg.ExecOptions{StreamPrint: false})
	if restoreExec.Error != nil {
		log.Println("Failed to restore")
		log.Println(restoreExec.Error.Err)
		log.Println(restoreExec.Output)

	} else {
		log.Println(restoreExec.Output)
		log.Println("Restore success")

	}
	err = spinner.Stop()
	if err != nil {
		log.Println("Failed to stop spinner")
	}
}
