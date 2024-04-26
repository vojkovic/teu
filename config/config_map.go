package config

// Seperate each application into a map of name to application
func (c *TeuConfig) ApplicationMap() map[string]struct {
	Deploy  string
	Secrets []string
} {
	appMap := make(map[string]struct {
		Deploy  string
		Secrets []string
	})
	for _, app := range c.Applications {
		appMap[app.Name] = struct {
			Deploy  string
			Secrets []string
		}{
			Deploy:  app.Deploy,
			Secrets: app.Secrets,
		}
	}
	return appMap
}

// Seperate each secret into a map of name to secret
func (c *TeuConfig) SecretMap() map[string]string {
	secretMap := make(map[string]string)
	for _, app := range c.Applications {
		for _, secret := range app.Secrets {
			secretMap[secret] = app.Name
		}
	}
	return secretMap
}