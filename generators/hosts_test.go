package generators

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/battlecrate/distrobuilder/image"
	"github.com/battlecrate/distrobuilder/shared"
)

func TestHostsGeneratorRunLXC(t *testing.T) {
	cacheDir := filepath.Join(os.TempDir(), "distrobuilder-test")
	rootfsDir := filepath.Join(cacheDir, "rootfs")

	setup(t, cacheDir)
	defer teardown(cacheDir)

	generator := Get("hosts")
	require.Equal(t, HostsGenerator{}, generator)

	definition := shared.Definition{
		Image: shared.DefinitionImage{
			Distribution: "ubuntu",
			Release:      "artful",
		},
	}

	image := image.NewLXCImage(cacheDir, "", cacheDir, definition)

	err := os.MkdirAll(filepath.Join(cacheDir, "rootfs", "etc"), 0755)
	require.NoError(t, err)

	createTestFile(t, filepath.Join(cacheDir, "rootfs", "etc", "hosts"),
		"127.0.0.1\tlocalhost\n127.0.0.1\tdistrobuilder\n")

	err = generator.RunLXC(cacheDir, rootfsDir, image, shared.DefinitionTargetLXC{},
		shared.DefinitionFile{Path: "/etc/hosts"})
	require.NoError(t, err)

	validateTestFile(t, filepath.Join(cacheDir, "rootfs", "etc", "hosts"),
		"127.0.0.1\tlocalhost\n127.0.0.1\tLXC_NAME\n")
}

func TestHostsGeneratorRunLXD(t *testing.T) {
	cacheDir := filepath.Join(os.TempDir(), "distrobuilder-test")
	rootfsDir := filepath.Join(cacheDir, "rootfs")

	setup(t, cacheDir)
	defer teardown(cacheDir)

	generator := Get("hosts")
	require.Equal(t, HostsGenerator{}, generator)

	definition := shared.Definition{
		Image: shared.DefinitionImage{
			Distribution: "ubuntu",
			Release:      "artful",
		},
	}

	image := image.NewLXDImage(cacheDir, "", cacheDir, definition)

	err := os.MkdirAll(filepath.Join(cacheDir, "rootfs", "etc"), 0755)
	require.NoError(t, err)

	createTestFile(t, filepath.Join(cacheDir, "rootfs", "etc", "hosts"),
		"127.0.0.1\tlocalhost\n127.0.0.1\tdistrobuilder\n")

	err = generator.RunLXD(cacheDir, rootfsDir, image, shared.DefinitionTargetLXD{},
		shared.DefinitionFile{Path: "/etc/hosts"})
	require.NoError(t, err)

	validateTestFile(t, filepath.Join(cacheDir, "templates", "hosts.tpl"),
		"127.0.0.1\tlocalhost\n127.0.0.1\t{{ container.name }}\n")
}
