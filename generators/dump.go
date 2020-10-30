package generators

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/flosch/pongo2"

	"github.com/battlecrate/distrobuilder/image"
	"github.com/battlecrate/distrobuilder/shared"
)

// DumpGenerator represents the Remove generator.
type DumpGenerator struct{}

// RunLXC dumps content to a file.
func (g DumpGenerator) RunLXC(cacheDir, sourceDir string, img *image.LXCImage,
	target shared.DefinitionTargetLXC, defFile shared.DefinitionFile) error {
	content := defFile.Content

	if defFile.Pongo {
		tpl, err := pongo2.FromString(defFile.Content)
		if err != nil {
			return err
		}

		content, err = tpl.Execute(pongo2.Context{"lxc": target})
		if err != nil {
			return err
		}
	}

	err := g.run(cacheDir, sourceDir, defFile, content)
	if err != nil {
		return err
	}

	if defFile.Templated {
		return img.AddTemplate(defFile.Path)
	}

	return nil
}

// RunLXD dumps content to a file.
func (g DumpGenerator) RunLXD(cacheDir, sourceDir string, img *image.LXDImage,
	target shared.DefinitionTargetLXD, defFile shared.DefinitionFile) error {
	content := defFile.Content

	if defFile.Pongo {
		tpl, err := pongo2.FromString(defFile.Content)
		if err != nil {
			return err
		}

		content, err = tpl.Execute(pongo2.Context{"lxd": target})
		if err != nil {
			return err
		}
	}

	return g.run(cacheDir, sourceDir, defFile, content)
}

// Run dumps content to a file.
func (g DumpGenerator) Run(cacheDir, sourceDir string, defFile shared.DefinitionFile) error {
	return g.run(cacheDir, sourceDir, defFile, defFile.Content)
}

func (g DumpGenerator) run(cacheDir, sourceDir string, defFile shared.DefinitionFile, content string) error {
	path := filepath.Join(sourceDir, defFile.Path)

	// Create any missing directory
	err := os.MkdirAll(filepath.Dir(path), 0755)
	if err != nil {
		return err
	}

	// Open the target file (create if needed)
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// Append final new line if missing
	if !strings.HasSuffix(content, "\n") {
		content += "\n"
	}

	// Write the content
	_, err = file.WriteString(content)
	if err != nil {
		return err
	}

	return updateFileAccess(file, defFile)
}
