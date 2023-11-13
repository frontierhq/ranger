package core

import (
	"regexp"
	"strings"
)

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

func (c *SetCollection) hasSet(s *Set) bool {
	for _, v := range c.Sets {
		if v.Name == s.Name && v.Environment == s.Environment {
			return true
		}
	}
	return false
}

func (c *SetCollection) addSet(s *Set) {
	if c.hasSet(s) {
		return
	}
	c.Sets = append(c.Sets, s)
}

func getSetNameFromRepoName(repoName *string) *string {
	n := strings.ReplaceAll(*repoName, "-set", "")
	re := regexp.MustCompile(`^.+?-`)
	n = re.ReplaceAllString(n, "")

	return &n
}
