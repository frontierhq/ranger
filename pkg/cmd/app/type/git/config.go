package git

// SetConfig sets the value for the given key
func (g *Git) SetConfig(key string, value string) error {
	_, err := g.Exec("config", "--local", key, value)
	if err != nil {
		return err
	}

	return nil
}
