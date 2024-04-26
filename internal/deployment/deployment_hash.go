package deployment

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"io"
	"os"
	"path/filepath"
)

// GetApplicationHash returns a hash of the entire deploy folder
func GetApplicationHash(folderPath string) (string, error) {
	hash := md5.New()
	isEmpty := true

	err := filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
			if info == nil {
					return err
			}
			
			if !info.IsDir() {
					isEmpty = false
					
					file, err := os.Open(path)
					if err != nil {
							return err
					}
					defer file.Close()

					if _, err := io.Copy(hash, file); err != nil {
							return err
					}
			}

			return nil
	})

	if err != nil {
			return "", err
	}

	if isEmpty {
			return "", errors.New("folder is empty")
	}

	hashInBytes := hash.Sum(nil)
	hashString := hex.EncodeToString(hashInBytes)

	return hashString, nil
}
