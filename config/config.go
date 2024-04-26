package config

type TeuConfig struct {
	Teu struct {
		Name          string `yaml:"name"`
		Description   string `yaml:"description"`
		AgeSecretKey  string `yaml:"age_secret_key"`
	} `yaml:"teu"`
	Applications []struct {
		Name    string   `yaml:"name"`
		Deploy  string   `yaml:"deploy"`
		Secrets []string `yaml:"secrets"`
	} `yaml:"applications"`
}
