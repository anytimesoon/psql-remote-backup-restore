package main

import (
	"flag"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
	"github.com/theckman/yacspin"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"sort"
)

var (
	local  db
	remote db

	pathToDump    string
	pathToRestore string
	pathToBackups string

	shouldBackup  bool
	shouldRestore bool

	backupOptions  []string
	restoreOptions []string

	shouldCleanUp bool
	maxBackups    int

	spinner *yacspin.Spinner

	restoreFile string
)

func init() {
	configName := flag.String("config", "main", "custom config name")
	flag.Parse()
	ex, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	dir := path.Dir(ex)
	log.Print(dir)

	viper.SetConfigName(*configName)
	viper.AddConfigPath(dir)
	viper.AddConfigPath("./")
	err = viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Fatalf("Could not find config file. Please make sure it is in the same directory as the executable and named %s.yaml.\n %s", *configName, err)
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
	pathToBackups = viper.GetString("directories.backups")
	if pathToBackups == "" {
		pathToBackups = dir + "backups/"
	}
	if pathToBackups[len(pathToBackups)-1:] != "/" {
		pathToBackups = pathToBackups + "/"
	}

	backupOptions = viper.GetStringSlice("backupOptions")
	restoreOptions = viper.GetStringSlice("restoreOptions")

	viper.SetDefault("shouldRestore", false)
	viper.SetDefault("shouldBackup", true)
	viper.SetDefault("shouldCleanUp", true)
	viper.SetDefault("maxBackups", 10)

	shouldRestore = viper.GetBool("shouldRestore")
	shouldBackup = viper.GetBool("shouldBackup")

	args := flag.Args()
	for _, arg := range args {
		switch arg {
		case "r", "restore":
			shouldBackup = false
		case "b", "backup":
			shouldRestore = false
		}
	}

	shouldCleanUp = viper.GetBool("shouldCleanUp")
	maxBackups = viper.GetInt("maxBackups")

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
			log.Println("No backup file explicitly given, searching for most recent backup")
			files, err := os.ReadDir(pathToBackups)
			if err != nil {
				log.Fatal(err)
			}
			sort.Slice(files, func(i, j int) bool {
				return files[i].Name() > files[j].Name()
			})
			restoreFile = pathToBackups + files[0].Name()
			log.Printf("Latest backup found was: %s", files[0].Name())
		}
		restore(restoreFile)
	}

	if shouldCleanUp {
		cleanUp()
	}
}

func findExecutables() {
	err := filepath.WalkDir(viper.GetString("directories.executables"), func(path string, d fs.DirEntry, err error) error {
		switch d.Name() {
		case "pg_dump":
			pathToDump = path
		case "pg_restore":
			pathToRestore = path
		default:

		}

		return nil
	})

	if err != nil || pathToRestore == "" || pathToDump == "" {
		log.Println("Failed to find the postgres executables (pg_dump and/or pg_restore)")
		log.Fatal(err)
	}
}
