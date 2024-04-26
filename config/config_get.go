package config

func (c *TeuConfig) GetSecrets() []string {
	var secrets []string
	for _, app := range c.Applications {
		secrets = append(secrets, app.Secrets...)
	}
	return secrets
}

func (c *TeuConfig) GetApplications() []string {
	var applications []string
	for _, app := range c.Applications {
		applications = append(applications, app.Name)
	}
	return applications
}

