package main

import "os"

func CreateFolder(name string) error {
	// Check if the destination folder exists, create it if not
	if _, err := os.Stat(name); os.IsNotExist(err) {
		if err := os.Mkdir(name, 0755); err != nil {
			return err
		}
	}

	return nil
}
