package git

// Checkout checks out a branch
func (g *Git) Checkout(branchName string, create bool) error {
	var args []string
	if create {
		args = []string{"checkout", "-b", branchName}
	} else {
		args = []string{"checkout", branchName}
	}

	_, err := g.Exec(args...)
	if err != nil {
		return err
	}

	return nil
}
