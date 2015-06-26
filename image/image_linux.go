package image

import (
	"encoding/json"
	"os"

	"github.com/docker/docker/pkg/archive"
)

// StoreImage stores file system layer data for the given image to the
// image's registered storage driver. Image metadata is stored in a file
// at the specified root directory.
func StoreImage(img *Image, layerData archive.ArchiveReader, root string) (err error) {
	// Store the layer. If layerData is not nil, unpack it into the new layer
	if layerData != nil {
		if img.Size, err = img.graph.Driver().ApplyDiff(img.ID, img.Parent, layerData); err != nil {
			return err
		}
	}

	if err := img.SaveSize(root); err != nil {
		return err
	}

	f, err := os.OpenFile(jsonPath(root), os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.FileMode(0600))
	if err != nil {
		return err
	}

	defer f.Close()

	return json.NewEncoder(f).Encode(img)
}

// TarLayer returns a tar archive of the image's filesystem layer.
func (img *Image) TarLayer() (arch archive.Archive, err error) {
	if img.graph == nil {
		return nil, fmt.Errorf("Can't load storage driver for unregistered image %s", img.ID)
	}

	driver := img.graph.Driver()

	return driver.Diff(img.ID, img.Parent)
}
