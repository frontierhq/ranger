package core

func (c *SetCollection) addSet(s *Set) {
	if c.hasSet(s) {
		return
	}
	c.Sets = append(c.Sets, s)
}

func (c *SetCollection) hasSet(s *Set) bool {
	for _, v := range c.Sets {
		if v.Name == s.Name && v.Environment == s.Environment {
			return true
		}
	}
	return false
}

func newSetCollection(sets *[]SetCollection, name string) []SetCollection {
	s := append(*sets, SetCollection{Name: name})
	return s
}

func getSetCollectionByName(sets *[]SetCollection, name string) *SetCollection {
	for _, s := range *sets {
		if s.Name == name {
			return &s
		}
	}
	return nil
}

func (c *SetCollection) findLastEnv() *Set {
	for _, s := range c.Sets {
		if s.Manifest.NextEnvironment == "" {
			return s
		}
	}
	return nil
}

func (c *SetCollection) findPreviousEnv(s *Set) *Set {
	for _, v := range c.Sets {
		if v.Environment != s.Environment && v.Manifest.NextEnvironment == s.Environment {
			return v
		}
	}
	return nil
}

func (c *SetCollection) Environments() []string {
	var envs []string
	current := c.Entry
	for {
		envs = append(envs, current.Environment)
		if current.Next == nil {
			break
		}
		current = current.Next
	}
	return envs
}

func (c *SetCollection) Order() {
	currentSet := c.findLastEnv()
	for {
		previousSet := c.findPreviousEnv(currentSet)
		if previousSet == nil {
			break
		}
		currentSet.Previous = previousSet
		previousSet.Next = currentSet
		currentSet = previousSet
	}
	c.Entry = currentSet
}
