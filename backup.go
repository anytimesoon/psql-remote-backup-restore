package main

import (
	pg "github.com/habx/pg-commands"
	"log"
	"os"
	"time"
)

func backup() string {
	pg.PGDumpCmd = pathToDump
	dump, err := pg.NewDump(&pg.Postgres{
		Host:     remote.host,
		Port:     remote.port,
		DB:       remote.name,
		Username: remote.user,
		Password: remote.password,
	})

	if err != nil {
		log.Fatal(err)
	}

	dump.Options = backupOptions
	dump.SetFileName("backups/" + time.Now().Format("20060102_1504") + "_backup" + ".sql")

	err = os.MkdirAll("backups", 0755)
	if err != nil {
		log.Println("Failed to create backups folder")
		log.Fatal(err)
	}

	log.Println("Starting backup")

	err = spinner.Start()
	spinner.Message("downloading data")

	dumpExec := dump.Exec(pg.ExecOptions{StreamPrint: false})
	if dumpExec.Error != nil {
		log.Println("Failed to dump")
		log.Println(dumpExec.Error.Err)
		log.Fatalln(dumpExec.Output)

	} else {
		log.Println(dumpExec.Output)
		log.Println("Backup success")
	}

	err = spinner.Stop()

	return dumpExec.File
}
