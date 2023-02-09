package git

// Add stages untracked files
func (g *Git) Add(path string) error {
	_, err := g.Exec("add", path)
	if err != nil {
		return err
	}

	return nil
}

// AddAll stages all untracked files
func (g *Git) AddAll() error {
	_, err := g.Exec("add", "--all")
	if err != nil {
		return err
	}

	return nil
}
