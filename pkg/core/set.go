package core

import (
	"regexp"
	"strings"
)

func getSetNameFromRepoName(repoName *string) *string {
	n := strings.ReplaceAll(*repoName, "-set", "")
	re := regexp.MustCompile(`^.+?-`)
	n = re.ReplaceAllString(n, "")

	return &n
}

func (s *Set) GetEnvNames() []string {
	var envs []string
	envs = append(envs, s.Environment)

	for {
		s = s.Next
		envs = append(envs, s.Environment)
		if s.Next == nil {
			break
		}
	}

	return envs
}
