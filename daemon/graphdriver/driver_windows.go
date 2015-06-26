// +build windows

package graphdriver

import (
	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/pkg/hcsshim"
)

type WindowsGraphDriver interface {
	Driver
	CopyDiff(id, sourceId string, parentLayerPaths []string) error
	LayerIdsToPaths(ids []string) []string
	Info() hcsshim.DriverInfo
	Export(id string, parentLayerPaths []string) (archive.Archive, error)
	Import(id string, layerData archive.ArchiveReader, parentLayerPaths []string) (int64, error)
}

const (
	FsMagicWindows       = FsMagic(0xa1b1830f) // I have just made this up for now. NTFS=0x5346544E
	FsMagicWindowsFilter = FsMagic(0xa1b1831f) // I have just made this up for now. NTFS=0x5346544E
	FsMagicWindowsDummy  = FsMagic(0xa1b1832f) // I have just made this up for now. NTFS=0x5346544E
)

var (
	// Slice of drivers that should be used in an order
	priority = []string{
		"windowsfilter",
		"windows",
		"windowsdummy",
	}

	FsNames = map[FsMagic]string{

		FsMagicWindows:       "windows",
		FsMagicWindowsFilter: "windowsfilter",
		FsMagicWindowsDummy:  "windowsdummy",
		FsMagicUnsupported:   "unsupported",
	}
)

func GetFSMagic(rootpath string) (FsMagic, error) {
	log.Debugln("WindowsGraphDriver GetFSMagic()")
	// TODO Windows
	return 0, nil
}
