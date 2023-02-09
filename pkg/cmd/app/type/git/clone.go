package git

import (
	netUrl "net/url"
)

// CloneOverHttp clones a Git repository using HTTP(S)
func (g *Git) CloneOverHttp(url string, username string, password string) error {
	parsedUrl, err := netUrl.Parse(url)
	if err != nil {
		return err
	}

	parsedUrl.User = netUrl.UserPassword(username, password)

	// err = os.RemoveAll(g.RepositoryPath)
	// if err != nil {
	// 	return err
	// }

	_, err = g.Exec("clone", parsedUrl.String(), g.repositoryPath)
	if err != nil {
		return err
	}

	return nil
}
