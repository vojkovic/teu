package git

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/vojkovic/teu/internal/db"
)

// Clone clones the repository if it doesn't exist, otherwise fetches the latest changes
func Update(repo, token string) error {
	// set last pull time in unix seconds
	err := db.SetRepositoryLastPullInDatabase(fmt.Sprintf("%d", time.Now().Unix()))
	if err != nil {
			return fmt.Errorf("error setting last pull time in database: %s", err)
	}

	fs := filepath.Join(os.Getenv("HOME"), ".teu", "repo")

	// Check if the repository already exists
	if _, err := os.Stat(fs); os.IsNotExist(err) {
			// Repository doesn't exist, clone it
			_, err := git.PlainClone(fs, false, &git.CloneOptions{
					URL: repo,
					Auth: &http.BasicAuth{
							Username: "username",
							Password: token,
					},
					Progress: os.Stdout,
					Depth:    1,
			})
			if err != nil {
					return fmt.Errorf("error cloning repository: %s", err)
			}
			return nil
	}

	// Open the existing repository
	r, err := git.PlainOpen(fs)
	if err != nil {
			return fmt.Errorf("error opening repository: %s", err)
	}

	// Get the current commit hash
	head, err := r.Head()
	if err != nil {
			return fmt.Errorf("error getting HEAD reference: %s", err)
	}
	currentHash := head.Hash()

	// Fetch the latest changes
	err = r.Fetch(&git.FetchOptions{
			Auth: &http.BasicAuth{
					Username: "username",
					Password: token,
			},
			Progress: os.Stdout,
	})
	if err != nil && err != git.NoErrAlreadyUpToDate {
			return fmt.Errorf("error fetching latest changes: %s", err)
	}

	// Get the new commit hash after fetch
	head, err = r.Head()
	if err != nil {
			return fmt.Errorf("error getting HEAD reference after fetch: %s", err)
	}
	newHash := head.Hash()


	fmt.Print("Current commit hash: " + currentHash.String() + "\n")
	fmt.Print("New commit hash: " + newHash.String() + "\n")
	// Compare the current hash with the new hash
	if currentHash.String() != newHash.String() {
			// If the hashes are different, remove the entire repository directory and reclone
			if err := os.RemoveAll(fs); err != nil {
					return fmt.Errorf("error removing repository directory: %s", err)
			}
			err := db.SetRepositoryLastCommitInDatabase(fmt.Sprintf("%d", time.Now().Unix()))
			if err != nil {
					return fmt.Errorf("error setting last commit time in database: %s", err)
			}
			// Clone the repository again
			_, err = git.PlainClone(fs, false, &git.CloneOptions{
					URL: repo,
					Auth: &http.BasicAuth{
							Username: "username",
							Password: token,
					},
					Progress: os.Stdout,
					Depth:    1,
			})
			if err != nil {
					return fmt.Errorf("error re-cloning repository: %s", err)
			}
	} else {
			fmt.Println("Repository is already up to date.")
	}

	return nil
}
