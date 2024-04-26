package config

import (
	"fmt"
	"os"
	"path/filepath"
)

func (c *TeuConfig) Validate(path string) error {
	if err := c.ValidateTeu(path); err != nil {
		return err
	}
	if err := c.ValidateSecret(path); err != nil {
		return err
	}
	return nil
}

func (c *TeuConfig) ValidateSecret(path string) error {
	if len(c.Teu.AgeSecretKey) == 0 {
		for _, app := range c.Applications {
			if len(app.Secrets) > 0 {
				return fmt.Errorf("age_secret_key is required in the config file")
			}
		}
	}
	// Check that each secret path is a valid file
	for _, app := range c.Applications {
		for _, secret := range app.Secrets {
			// /Users/courage/Sync/Personal/Teu/sl8/secret.key
			completePath := filepath.Join(path, app.Deploy, secret)
			if _, err := os.Stat(completePath); err != nil {
				return fmt.Errorf("secret path %s is not a valid file", completePath)
			}
		}
	}

	return nil
}

func (c *TeuConfig) ValidateTeu(path string) error {
	if len(c.Teu.Name) == 0 {
		return fmt.Errorf("name is required in the config file")
	}
	if len(c.Teu.Description) == 0 {
		return fmt.Errorf("description is required in the config file")
	}
	if (len(c.Applications) == 0) {
		return fmt.Errorf("at least 1 application is required in the config file")
	}
	return nil
}




