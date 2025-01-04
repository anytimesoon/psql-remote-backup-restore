package main

import (
	"log"
	"os"
	"path/filepath"
	"sort"
)

func cleanUp() {
	files, err := os.ReadDir(pathToBackups)
	if err != nil {
		log.Println("Failed to read backup directory for clean up")
		return
	}

	if len(files) > maxBackups {
		sort.Slice(files, func(i, j int) bool {
			info1, err1 := files[i].Info()
			info2, err2 := files[j].Info()
			if err1 != nil || err2 != nil {
				log.Printf("Error retrieving file info: %v, %v", err1, err2)
				return false
			}

			return info1.ModTime().Before(info2.ModTime())
		})

		// The oldest file will now be the first element in the sorted list
		oldestFile := files[0]
		oldestFilePath := filepath.Join(pathToBackups, oldestFile.Name())

		// Delete the oldest file
		log.Printf("Deleting oldest file: %s\n", oldestFilePath)
		err = os.Remove(oldestFilePath)
		if err != nil {
			log.Fatalf("Failed to delete file %s: %v", oldestFilePath, err)
		}

		log.Printf("Oldest file %s has been deleted successfully.\n", oldestFile.Name())
		return
	}

	log.Println("Not enough backups to require clean up")
}
