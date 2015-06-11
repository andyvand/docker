package graph

import (
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/daemon/graphdriver/windows"
	"github.com/docker/docker/image"
	"github.com/docker/docker/pkg/archive"
	"github.com/docker/docker/utils"
)

// setupInitLayer populates a directory with mountpoints suitable
// for bind-mounting dockerinit into the container. The mountpoint is simply an
// empty file at /.dockerinit
//
// This extra layer is used by all containers as the top-most ro layer. It protects
// the container from unwanted side-effects on the rw layer.
func SetupInitLayer(initLayer string) error {
	log.Debugln("Windows SetupInitLayer()")
	return nil
}

// Register imports a pre-existing image into the graph.
func (graph *Graph) Register(img *image.Image, layerData archive.ArchiveReader) (err error) {
	defer func() {
		// If any error occurs, remove the new dir from the driver.
		// Don't check for errors since the dir might not have been created.
		// FIXME: this leaves a possible race condition.
		if err != nil {
			graph.driver.Remove(img.ID)
		}
	}()
	if err := utils.ValidateID(img.ID); err != nil {
		return err
	}
	// (This is a convenience to save time. Race conditions are taken care of by os.Rename)
	if graph.Exists(img.ID) {
		return fmt.Errorf("Image %s already exists", img.ID)
	}

	// Ensure that the image root does not exist on the filesystem
	// when it is not registered in the graph.
	// This is common when you switch from one graph driver to another
	if err := os.RemoveAll(graph.ImageRoot(img.ID)); err != nil && !os.IsNotExist(err) {
		return err
	}

	// If the driver has this ID but the graph doesn't, remove it from the driver to start fresh.
	// (the graph is the source of truth).
	// Ignore errors, since we don't know if the driver correctly returns ErrNotExist.
	// (FIXME: make that mandatory for drivers).
	graph.driver.Remove(img.ID)

	tmp, err := graph.Mktemp("")
	defer os.RemoveAll(tmp)
	if err != nil {
		return fmt.Errorf("Mktemp failed: %s", err)
	}

	if wd, ok := graph.driver.(*windows.WindowsGraphDriver); ok && img.Container != "" && layerData == nil {
		log.Debugf("Copying from container %s.", img.Container)

		var ids []string
		if img.Parent != "" {
			parentImg, err := graph.Get(img.Parent)
			if err != nil {
				return err
			}

			ids, err = graph.ParentLayerIds(parentImg)
			if err != nil {
				return err
			}
		}

		if err := wd.CopyDiff(img.Container, img.ID, wd.LayerIdsToPaths(ids)); err != nil {
			return fmt.Errorf("Driver %s failed to copy image rootfs %s: %s", graph.driver, img.Container, err)
		}
	} else {
		// Create root filesystem in the driver
		if err := graph.driver.Create(img.ID, img.Parent); err != nil {
			return fmt.Errorf("Driver %s failed to create image rootfs %s: %s", graph.driver, img.ID, err)
		}
	}

	// Apply the diff/layer
	img.SetGraph(graph)
	if err := image.StoreImage(img, layerData, tmp); err != nil {
		return err
	}

	// Commit
	if err := os.Rename(tmp, graph.ImageRoot(img.ID)); err != nil {
		return err
	}
	graph.idIndex.Add(img.ID)
	return nil
}

func (graph *Graph) restoreBaseImages() ([]string, error) {
	return nil, nil
}
