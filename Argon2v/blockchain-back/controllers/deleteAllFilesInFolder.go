package controllers

import (
	"os"
	"path/filepath"
)

func DeleteAllFilesInFolder(folderPath string) error {
	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			err = os.Remove(path)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return err
}
