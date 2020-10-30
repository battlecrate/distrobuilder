package managers

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	lxd "github.com/lxc/lxd/shared"

	"github.com/battlecrate/distrobuilder/shared"
)

// NewYum creates a new Manager instance.
func NewYum() *Manager {
	var buf bytes.Buffer
	globalFlags := []string{"-y"}

	lxd.RunCommandWithFds(nil, &buf, "yum", "--help")

	scanner := bufio.NewScanner(&buf)

	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "--allowerasing") {
			globalFlags = append(globalFlags, "--allowerasing")
			break
		}
	}

	return &Manager{
		commands: ManagerCommands{
			clean:   "yum",
			install: "yum",
			refresh: "yum",
			remove:  "yum",
			update:  "yum",
		},
		flags: ManagerFlags{
			clean: []string{
				"clean", "all",
			},
			global: globalFlags,
			install: []string{
				"install",
			},
			remove: []string{
				"remove",
			},
			refresh: []string{
				"makecache",
			},
			update: []string{
				"update",
			},
		},
		RepoHandler: yumRepoHandler,
	}
}

func yumRepoHandler(repoAction shared.DefinitionPackagesRepository) error {
	targetFile := filepath.Join("/etc/yum.repos.d", repoAction.Name)

	if !strings.HasSuffix(targetFile, ".repo") {
		targetFile = fmt.Sprintf("%s.repo", targetFile)
	}

	if !lxd.PathExists(filepath.Dir(targetFile)) {
		err := os.MkdirAll(filepath.Dir(targetFile), 0755)
		if err != nil {
			return err
		}
	}

	f, err := os.Create(targetFile)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(repoAction.URL)
	if err != nil {
		return err
	}

	// Append final new line if missing
	if !strings.HasSuffix(repoAction.URL, "\n") {
		_, err = f.WriteString("\n")
		if err != nil {
			return err
		}
	}

	return nil
}
