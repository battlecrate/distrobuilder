package generators

import (
	"fmt"
	"github.com/battlecrate/distrobuilder/shared"
	"io"
	"os"
	"path/filepath"

	"github.com/battlecrate/distrobuilder/image"
)

// CopyGenerator represents the Copy generator.
type CopyGenerator struct{}

// RunLXC copies a file to the container.
func (g CopyGenerator) RunLXC(cacheDir string, sourceDir string, img *image.LXCImage,
	target shared.DefinitionTargetLXC, defFile shared.DefinitionFile) error {
	// no template support for LXC, ignoring generator
	return nil
}

// RunLXD copies a file to the container.
func (g CopyGenerator) RunLXD(cacheDir, sourceDir string, img *image.LXDImage,
	target shared.DefinitionTargetLXD, defFile shared.DefinitionFile) error {
	return g.Run(cacheDir, sourceDir, defFile)
}

// Run copies a file to the container.
func (g CopyGenerator) Run(cacheDir, sourceDir string, defFile shared.DefinitionFile) error {
	in, err := os.Open(defFile.Source)
	if err != nil {
		if os.IsNotExist(err) {
			err = fmt.Errorf("File '%s' doesn't exist", defFile.Path)
		}
		return err
	}
	defer in.Close()

	// Calculate the destination path by combining the "rootfs" path and the destination file
	var path = filepath.Join(sourceDir, defFile.Path)

	// Create any missing directory
	err = os.MkdirAll(filepath.Dir(path), 0755)
	if err != nil {
		return err
	}

	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()

	// copy data onto the new file
	_, err = io.Copy(out, in)

	if err != nil {
		return err
	}

	// update the file access permissions
	err = updateFileAccess(out, defFile)
	return err
}
