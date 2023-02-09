package git

// Commit commits all staged changes
func (g *Git) Commit(message string) error {
	_, err := g.Exec("commit", "-m", message)
	if err != nil {
		return err
	}

	return nil
}
