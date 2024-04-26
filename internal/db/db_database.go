package db

type Database struct {
	Repository   					string           				`json:"repository"`
	RepositoryToken   		string            			`json:"token"`
	RepositoryLastPull  	string            			`json:"last_pull"`
	RepositoryLastCommit 	string            			`json:"last_commit"`
	Applications 					map[string]Application 	`json:"applications"`
}

type Application struct {
	Hash   string            `json:"hash"`
}
