package git

// Push pushes the current branch to the "origin" remote
func (g *Git) Push(force bool) error {
	var args []string
	if force {
		args = []string{"push", "-u", "origin", "HEAD"}
	} else {
		args = []string{"push", "-u", "origin", "HEAD", "--force"}
	}

	_, err := g.Exec(args...)
	if err != nil {
		return err
	}

	return nil
}
