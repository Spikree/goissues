package helper

import (
	"os/exec"
	"strings"
)

// gitUser returns the best available identity from git config.
func GitUser() string {
	for _, key := range []string{"github.user", "user.name"} {
		out, err := exec.Command("git", "config", "--get", key).Output()
		if err == nil {
			if name := strings.TrimSpace(string(out)); name != "" {
				return name
			}
		}
	}
	return "unknown"
}
