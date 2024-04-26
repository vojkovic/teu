package age

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// DecryptSecret decrypts the secrets in the config file
// using the age secret key
// age --decrypt --output ./secret.key ./secret.key.age
// Output is just input -.age
func DecryptSecret(secret_path, path string)(error) {

	err := isAgeSecretKey(secret_path)
	if err != nil {
		return err
	}

	output_path, err := removeAgeSuffix(path)
	if err != nil {
		return err
	}

	err = isAgeEncryptedFile(path)
	if err != nil {
		return err
	}

	cmd := exec.Command("age", "--decrypt", "-i", secret_path, "--output", output_path, path)
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("could not decrypt the secret: %w", err)
	}

	return nil
}

func removeAgeSuffix(path string) (string, error) {
	if strings.HasSuffix(path, ".age") {
		return path[:len(path)-4], nil
	} else {
		return "", fmt.Errorf("path does not have the .age suffix, %s", path)
	}
}


// verify the key is a valid age secret key by checking if the file has a single: AGE-SECRET-KEY within the file somewhere.
func isAgeSecretKey(path string) (error) {
	info, err := os.Stat(filepath.Clean(path))
	if err != nil {
		return fmt.Errorf("could not find the age key in %s: %s", path, err)
	}

	if info.IsDir() {
		return fmt.Errorf("file %s is a directory, expected file", path)
	}
	
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("could not open the file of the age secret key: %s", path)
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return fmt.Errorf("could not get the file stats of the age secret key: %s", path)
	}

	key := make([]byte, stat.Size())
	_, err = file.Read(key)
	if err != nil {
		return fmt.Errorf("could not read the key from the file containing the age secret key: %s", path)
	}

	if !strings.Contains(string(key), "AGE-SECRET-KEY") {
		return fmt.Errorf("file does not contain a valid age secret key: %s", path)
	}

	return nil
}

// verify the key is a valid age secret key by checking if the file has a ascii armored age encrypted file.
func isAgeEncryptedFile(path string)(error) {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("could not open the age encrypted file: %s", path)
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return fmt.Errorf("could not get the age encrypted file: %s", path)
	}

	data := make([]byte, stat.Size())
	_, err = file.Read(data)
	if err != nil {
		return fmt.Errorf("could not read the age encrypted file: %s", path)
	}

	if !strings.Contains(string(data), "-----BEGIN AGE ENCRYPTED FILE-----") {
		return fmt.Errorf("file is not age encrypted, no header: %s", path)
	}

	if !strings.Contains(string(data), "-----END AGE ENCRYPTED FILE-----") {
		return fmt.Errorf("file is not age encrypted, no footer: %s", path)
	}

	return nil
}