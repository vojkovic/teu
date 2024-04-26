package db

func GetRepositoryTokenFromDatabase() (string, error) {
	return SetKeyStringInDatabase("RepositoryToken")
}

func SetRepositoryTokenInDatabase(token string) error {
	return SetValueStringInDatabase("RepositoryToken", token)
}

func GetRepositoryFromDatabase() (string, error) {
	return SetKeyStringInDatabase("Repository")
}

func SetRepositoryInDatabase(repository string) error {
	return SetValueStringInDatabase("Repository", repository)
}

func GetRepositoryLastPullFromDatabase() (string, error) {
	return SetKeyStringInDatabase("RepositoryLastPull")
}

func SetRepositoryLastPullInDatabase(lastPull string) error {
	return SetValueStringInDatabase("RepositoryLastPull", lastPull)
}

func GetRepositoryLastCommitFromDatabase() (string, error) {
	return SetKeyStringInDatabase("RepositoryLastCommit")
}

func SetRepositoryLastCommitInDatabase(lastCommit string) error {
	return SetValueStringInDatabase("RepositoryLastCommit", lastCommit)
}
