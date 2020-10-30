package generators

import (
	"os"
	"path/filepath"

	"github.com/battlecrate/distrobuilder/image"
	"github.com/battlecrate/distrobuilder/shared"
)

// RemoveGenerator represents the Remove generator.
type RemoveGenerator struct{}

// RunLXC removes a path.
func (g RemoveGenerator) RunLXC(cacheDir, sourceDir string, img *image.LXCImage,
	target shared.DefinitionTargetLXC, defFile shared.DefinitionFile) error {
	return g.Run(cacheDir, sourceDir, defFile)
}

// RunLXD removes a path.
func (g RemoveGenerator) RunLXD(cacheDir, sourceDir string, img *image.LXDImage,
	target shared.DefinitionTargetLXD, defFile shared.DefinitionFile) error {
	return g.Run(cacheDir, sourceDir, defFile)
}

// Run removes a path.
func (g RemoveGenerator) Run(cacheDir, sourceDir string,
	defFile shared.DefinitionFile) error {
	return os.RemoveAll(filepath.Join(sourceDir, defFile.Path))
}
