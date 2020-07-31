package file

import "os"

// Exists A function to check of a file/folder exists
func Exists(filePath string) bool {
	_, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}

	return true
}
