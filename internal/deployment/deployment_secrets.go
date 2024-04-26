package deployment

import (
	"path/filepath"
	"strings"

	"github.com/vojkovic/teu/internal/age"
)

func DecryptSecretsInApp(secrets[] string, age_secret_key, deploy_location, app_name string) error {
	if len(secrets) == 0 {
		return nil
	}

	for _, secret := range secrets {
		enc_secret_name := strings.Split(secret, "/")[len(strings.Split(secret, "/"))-1]
		secret_name := enc_secret_name[:len(enc_secret_name)-4]

		PrintSecretDecrypting(enc_secret_name, secret_name)
		
		err := age.DecryptSecret(age_secret_key, filepath.Join(deploy_location, secret))
		if err != nil {
			PrintSecretFailedToDecrypt(enc_secret_name)
			return err
		}
	}
	return nil
}