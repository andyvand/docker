package archive

import (
	"github.com/docker/docker/pkg/system"
)

type FileInfo struct {
	parent     *FileInfo
	name       string
	stat       *system.Stat
	children   map[string]*FileInfo
	capability []byte
	added      bool
}

func statDifferent(oldStat *system.Stat, newStat *system.Stat) bool {

	// Don't look at size for dirs, its not a good measure of change
	if oldStat.ModTime() != newStat.ModTime() ||
		oldStat.Mode() != newStat.Mode() ||
		oldStat.Size() != newStat.Size() && !oldStat.IsDir() {
		return true
	} else {
		return false
	}
}

func (info *FileInfo) isDir() bool {
	return info.parent == nil || info.stat.IsDir()
}