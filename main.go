package main

import (
	"fmt"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"github.com/theckman/yacspin"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"sort"
)

var (
	local  db
	remote db

	pathToDump    string
	pathToRestore string

	shouldBackup  bool
	shouldRestore bool

	spinner *yacspin.Spinner

	restoreFile string
)

func init() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Fatalln("Could not find config file. Please make sure it is in the same directory as the executable and named config.yaml.", err)
		} else {
			log.Fatalln("Failed to read config file.", err)
		}
	}

	local = db{
		host:     viper.GetString("localDb.host"),
		port:     viper.GetInt("localDb.port"),
		user:     viper.GetString("localDb.user"),
		password: viper.GetString("localDb.password"),
		name:     viper.GetString("localDb.name"),
	}

	remote = db{
		host:     viper.GetString("remoteDb.host"),
		port:     viper.GetInt("remoteDb.port"),
		user:     viper.GetString("remoteDb.user"),
		password: viper.GetString("remoteDb.password"),
		name:     viper.GetString("remoteDb.name"),
	}

	findExecutables()

	shouldRestore = viper.GetBool("shouldRestore")
	shouldBackup = viper.GetBool("shouldBackup")

	spinner, err = yacspin.New(cfg)
	if err != nil {
		log.Println("Failed to start the spinner")
		log.Println(err)
	}
}

func main() {
	if shouldBackup {
		restoreFile = backup()
	}

	if shouldRestore {
		if restoreFile == "" {
			files, err := os.ReadDir("backups/")
			if err != nil {
				log.Fatal(err)
			}
			sort.Slice(files, func(i, j int) bool {
				return files[i].Name() > files[j].Name()
			})
			restoreFile = fmt.Sprintf("backups/%s", files[0].Name())
		}
		restore(restoreFile)
	}
}

func findExecutables() {
	err := filepath.WalkDir(viper.GetString("pathToExecutables"), func(path string, d fs.DirEntry, err error) error {
		switch d.Name() {
		case "pg_dump":
			pathToDump = filepath.Join(viper.GetString("pathToExecutables"), path, "pg_dump")
		case "pg_restore":
			pathToRestore = filepath.Join(viper.GetString("pathToExecutables"), path, "pg_dump")
		default:

		}

		return nil
	})

	if err != nil || pathToRestore == "" || pathToDump == "" {
		log.Println("Failed to find the postgres executables (pg_dump and/or pg_restore)")
		log.Fatal(err)
	}
}
