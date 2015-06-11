// +build windows

// Shim for Windows Containers Host Compute Service (HSC)
// This file contains the layer specific shim function calls, used to implement
// the Windows graphdriver functionality.

package hcsshim

import (
	"syscall"
	"unsafe"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/pkg/guid"
)

var (
	procLayerExists        = modvmcompute.NewProc("LayerExists")
	procCreateLayer        = modvmcompute.NewProc("CreateLayer")
	procDestroyLayer       = modvmcompute.NewProc("DestroyLayer")
	procActivateLayer      = modvmcompute.NewProc("ActivateLayer")
	procDeactivateLayer    = modvmcompute.NewProc("DeactivateLayer")
	procGetLayerMountPath  = modvmcompute.NewProc("GetLayerMountPath")
	procCopyLayer          = modvmcompute.NewProc("CopyLayer")
	procCreateSandboxLayer = modvmcompute.NewProc("CreateSandboxLayer")
	procPrepareLayer       = modvmcompute.NewProc("PrepareLayer")
	procUnprepareLayer     = modvmcompute.NewProc("UnprepareLayer")
)

/* To pass into syscall, we need a struct matching the following:
enum GraphDriverType
{
    DiffDriver,
    FilterDriver
};

struct DriverInfo {
    GraphDriverType Flavor;
    LPCWSTR HomeDir;
};
*/
type DriverInfo struct {
	Flavor  int
	HomeDir string
}

type driverInfo struct {
	Flavor   int
	HomeDirp *uint16
}

func convertInfo(info DriverInfo) (driverInfo, error) {
	homedirp, err := syscall.UTF16PtrFromString(info.HomeDir)
	if err != nil {
		log.Debugln("Failed conversion of home to pointer for driver info: ", err.Error())
		return driverInfo{}, err
	}

	return driverInfo{
		Flavor:   info.Flavor,
		HomeDirp: homedirp,
	}, nil
}

/* To pass into syscall, we need a struct matching the following:
typedef struct _WC_LAYER_DESCRIPTOR {

    //
    // The ID of the layer
    //

    GUID LayerId;

    //
    // Additional flags
    //

    union {
        struct {
            ULONG Reserved : 31;
            ULONG Dirty : 1;    // Created from sandbox as a result of snapshot
        };
        ULONG Value;
    } Flags;

    //
    // Path to the layer root directory, null-terminated
    //

    PCWSTR Path;

} WC_LAYER_DESCRIPTOR, *PWC_LAYER_DESCRIPTOR;
*/
type WC_LAYER_DESCRIPTOR struct {
	LayerId guid.Guid
	Flags   uint32
	Pathp   *uint16
}

func LayerPathsToDescriptors(parentLayerPaths []string) ([]WC_LAYER_DESCRIPTOR, error) {
	// Array of descriptors that gets constructed.
	var layers []WC_LAYER_DESCRIPTOR

	for i := 0; i < len(parentLayerPaths); i++ {
		// Create a layer descriptor, using the folder path
		// as the source for a GUID LayerId
		g := guid.NewGuid(parentLayerPaths[i])

		p, err := syscall.UTF16PtrFromString(parentLayerPaths[i])
		if err != nil {
			log.Debugln("Failed conversion of parentLayerPath to pointer ", err)
			return nil, err
		}

		layers = append(layers, WC_LAYER_DESCRIPTOR{
			LayerId: *g,
			Flags:   0,
			Pathp:   p,
		})
	}

	return layers, nil
}

func LayerExists(info DriverInfo, id string) (bool, error) {
	log.Debugln("hcsshim::LayerExists")
	log.Debugln("info.Flavor:", info.Flavor)
	log.Debugln("id:", id)

	idp, err := syscall.UTF16PtrFromString(id)
	if err != nil {
		log.Debugln("Failed conversion of id to pointer ", err)
		return false, err
	}

	infop, err := convertInfo(info)
	if err != nil {
		log.Debugln("Failed conversion of driver info ", err)
		return false, err
	}

	var exists bool

	// Call the procedure itself.
	r1, _, _ := procLayerExists.Call(
		uintptr(unsafe.Pointer(&infop)),
		uintptr(unsafe.Pointer(idp)),
		uintptr(unsafe.Pointer(&exists)))
	use(unsafe.Pointer(&infop))
	use(unsafe.Pointer(idp))

	if r1 != 0 {
		return false, syscall.Errno(r1)
	}

	return exists, nil
}

func CreateLayer(info DriverInfo, id, parent string) error {
	log.Debugln("hcsshim::CreateLayer")
	log.Debugln("info.Flavor:", info.Flavor)
	log.Debugln("id:", id)
	log.Debugln("parent:", parent)

	idp, err := syscall.UTF16PtrFromString(id)
	if err != nil {
		log.Debugln("Failed conversion of id to pointer ", err)
		return err
	}

	parentp, err := syscall.UTF16PtrFromString(parent)
	if err != nil {
		log.Debugln("Failed conversion of parent to pointer ", err)
		return err
	}

	infop, err := convertInfo(info)
	if err != nil {
		log.Debugln("Failed conversion of driver info ", err)
		return err
	}

	// Call the procedure itself.
	r1, _, _ := procCreateLayer.Call(
		uintptr(unsafe.Pointer(&infop)),
		uintptr(unsafe.Pointer(idp)),
		uintptr(unsafe.Pointer(parentp)))
	use(unsafe.Pointer(&infop))
	use(unsafe.Pointer(idp))
	use(unsafe.Pointer(parentp))

	if r1 != 0 {
		return syscall.Errno(r1)
	}

	return nil
}

func DestroyLayer(info DriverInfo, id string) error {
	log.Debugln("hcsshim::DestroyLayer")
	log.Debugln("info.Flavor:", info.Flavor)
	log.Debugln("id:", id)

	idp, err := syscall.UTF16PtrFromString(id)
	if err != nil {
		log.Debugln("Failed conversion of id to pointer ", err)
		return err
	}

	infop, err := convertInfo(info)
	if err != nil {
		log.Debugln("Failed conversion of driver info ", err)
		return err
	}

	// Call the procedure itself.
	r1, _, _ := procDestroyLayer.Call(
		uintptr(unsafe.Pointer(&infop)),
		uintptr(unsafe.Pointer(idp)))
	use(unsafe.Pointer(&infop))
	use(unsafe.Pointer(idp))

	if r1 != 0 {
		return syscall.Errno(r1)
	}

	return nil
}

func ActivateLayer(info DriverInfo, id string) error {
	log.Debugln("hcsshim::ActivateLayer")
	log.Debugln("info.Flavor:", info.Flavor)
	log.Debugln("id:", id)

	idp, err := syscall.UTF16PtrFromString(id)
	if err != nil {
		log.Debugln("Failed conversion of id to pointer ", err)
		return err
	}

	infop, err := convertInfo(info)
	if err != nil {
		log.Debugln("Failed conversion of driver info ", err)
		return err
	}

	// Call the procedure itself.
	r1, _, _ := procActivateLayer.Call(
		uintptr(unsafe.Pointer(&infop)),
		uintptr(unsafe.Pointer(idp)))
	use(unsafe.Pointer(&infop))
	use(unsafe.Pointer(idp))

	if r1 != 0 {
		return syscall.Errno(r1)
	}

	return nil
}

func DeactivateLayer(info DriverInfo, id string) error {
	log.Debugln("hcsshim::DeactivateLayer")
	log.Debugln("info.Flavor:", info.Flavor)
	log.Debugln("id:", id)

	idp, err := syscall.UTF16PtrFromString(id)
	if err != nil {
		log.Debugln("Failed conversion of id to pointer ", err)
		return err
	}

	infop, err := convertInfo(info)
	if err != nil {
		log.Debugln("Failed conversion of driver info ", err)
		return err
	}

	// Call the procedure itself.
	r1, _, _ := procDeactivateLayer.Call(
		uintptr(unsafe.Pointer(&infop)),
		uintptr(unsafe.Pointer(idp)))
	use(unsafe.Pointer(&infop))
	use(unsafe.Pointer(idp))

	if r1 != 0 {
		return syscall.Errno(r1)
	}

	return nil
}

func GetLayerMountPath(info DriverInfo, id string) (string, error) {
	log.Debugln("hcsshim::GetLayerMountPath")
	log.Debugln("info.Flavor:", info.Flavor)
	log.Debugln("id:", id)

	idp, err := syscall.UTF16PtrFromString(id)
	if err != nil {
		log.Debugln("Failed conversion of id to pointer ", err)
		return "", err
	}

	infop, err := convertInfo(info)
	if err != nil {
		log.Debugln("Failed conversion of driver info ", err)
		return "", err
	}

	var mountPathLength uint64
	mountPathLength = 0

	r1, _, _ := procGetLayerMountPath.Call(
		uintptr(unsafe.Pointer(&infop)),
		uintptr(unsafe.Pointer(idp)),
		uintptr(unsafe.Pointer(&mountPathLength)),
		uintptr(unsafe.Pointer(nil)))
	if r1 != 0 {
		return "", syscall.Errno(r1)
	}

	// Allocate a mount path of the returned length.
	if mountPathLength == 0 {
		return "", nil
	}

	mountPathp := make([]uint16, mountPathLength)
	mountPathp[0] = 0

	r1, _, _ = procGetLayerMountPath.Call(
		uintptr(unsafe.Pointer(&infop)),
		uintptr(unsafe.Pointer(idp)),
		uintptr(unsafe.Pointer(&mountPathLength)),
		uintptr(unsafe.Pointer(&mountPathp[0])))
	use(unsafe.Pointer(&mountPathLength))
	use(unsafe.Pointer(&infop))
	use(unsafe.Pointer(idp))

	if r1 != 0 {
		return "", syscall.Errno(r1)
	}

	return syscall.UTF16ToString(mountPathp[0:]), nil
}

func CopyLayer(info DriverInfo, srcId, dstId string, parentLayerPaths []string) error {
	log.Debugln("hcsshim::CopyLayer")
	log.Debugln("info.Flavor:", info.Flavor)
	log.Debugln("srcId:", srcId)
	log.Debugln("dstId:", dstId)

	layers, err := LayerPathsToDescriptors(parentLayerPaths)
	if err != nil {
		log.Debugln("Failed to generate layer descriptors ", err)
		return err
	}

	srcIdp, err := syscall.UTF16PtrFromString(srcId)
	if err != nil {
		log.Debugln("Failed conversion of srcId to pointer ", err)
		return err
	}

	dstIdp, err := syscall.UTF16PtrFromString(dstId)
	if err != nil {
		log.Debugln("Failed conversion of dstId to pointer ", err)
		return err
	}

	infop, err := convertInfo(info)
	if err != nil {
		log.Debugln("Failed conversion of driver info ", err)
		return err
	}

	var layerDescriptorsp *WC_LAYER_DESCRIPTOR
	if len(layers) > 0 {
		layerDescriptorsp = &(layers[0])
	} else {
		layerDescriptorsp = nil
	}

	// Call the procedure itself.
	r1, _, _ := procCopyLayer.Call(
		uintptr(unsafe.Pointer(&infop)),
		uintptr(unsafe.Pointer(srcIdp)),
		uintptr(unsafe.Pointer(dstIdp)),
		uintptr(unsafe.Pointer(layerDescriptorsp)),
		uintptr(len(layers)))
	use(unsafe.Pointer(&infop))
	use(unsafe.Pointer(srcIdp))
	use(unsafe.Pointer(dstIdp))
	use(unsafe.Pointer(layerDescriptorsp))

	if r1 != 0 {
		return syscall.Errno(r1)
	}

	return nil
}

func CreateSandboxLayer(info DriverInfo, layerId, parentId string, parentLayerPaths []string) error {
	log.Debugln("hcsshim::CreateSandboxLayer")
	log.Debugln("info.Flavor:", info.Flavor)
	log.Debugln("layerId:", layerId)

	layers, err := LayerPathsToDescriptors(parentLayerPaths)
	if err != nil {
		log.Debugln("Failed to generate layer descriptors ", err)
		return err
	}

	layerIdp, err := syscall.UTF16PtrFromString(layerId)
	if err != nil {
		log.Debugln("Failed conversion of layerId to pointer ", err)
		return err
	}

	parentIdp, err := syscall.UTF16PtrFromString(parentId)
	if err != nil {
		log.Debugln("Failed conversion of parentId to pointer ", err)
		return err
	}

	infop, err := convertInfo(info)
	if err != nil {
		log.Debugln("Failed conversion of driver info ", err)
		return err
	}

	var layerDescriptorsp *WC_LAYER_DESCRIPTOR
	if len(layers) > 0 {
		layerDescriptorsp = &(layers[0])
	} else {
		layerDescriptorsp = nil
	}

	// Call the procedure itself.
	r1, _, _ := procCreateSandboxLayer.Call(
		uintptr(unsafe.Pointer(&infop)),
		uintptr(unsafe.Pointer(layerIdp)),
		uintptr(unsafe.Pointer(parentIdp)),
		uintptr(unsafe.Pointer(layerDescriptorsp)),
		uintptr(len(layers)))
	use(unsafe.Pointer(&infop))
	use(unsafe.Pointer(layerIdp))
	use(unsafe.Pointer(parentIdp))
	use(unsafe.Pointer(layerDescriptorsp))

	if r1 != 0 {
		return syscall.Errno(r1)
	}

	return nil
}

func PrepareLayer(info DriverInfo, layerId string, parentLayerPaths []string) error {
	log.Debugln("hcsshim::PrepareLayer")
	log.Debugln("info.Flavor:", info.Flavor)
	log.Debugln("layerId:", layerId)

	layers, err := LayerPathsToDescriptors(parentLayerPaths)
	if err != nil {
		log.Debugln("Failed to generate layer descriptors ", err)
		return err
	}

	layerIdp, err := syscall.UTF16PtrFromString(layerId)
	if err != nil {
		log.Debugln("Failed conversion of layerId to pointer ", err)
		return err
	}

	infop, err := convertInfo(info)
	if err != nil {
		log.Debugln("Failed conversion of driver info ", err)
		return err
	}

	var layerDescriptorsp *WC_LAYER_DESCRIPTOR
	if len(layers) > 0 {
		layerDescriptorsp = &(layers[0])
	} else {
		layerDescriptorsp = nil
	}

	// Call the procedure itself.
	r1, _, _ := procPrepareLayer.Call(
		uintptr(unsafe.Pointer(&infop)),
		uintptr(unsafe.Pointer(layerIdp)),
		uintptr(unsafe.Pointer(layerDescriptorsp)),
		uintptr(len(layers)))
	use(unsafe.Pointer(&infop))
	use(unsafe.Pointer(layerIdp))
	use(unsafe.Pointer(layerDescriptorsp))

	if r1 != 0 {
		return syscall.Errno(r1)
	}

	return nil
}

func UnprepareLayer(info DriverInfo, layerId string) error {
	log.Debugln("hcsshim::UnprepareLayer")
	log.Debugln("info.Flavor:", info.Flavor)
	log.Debugln("layerId:", layerId)

	layerIdp, err := syscall.UTF16PtrFromString(layerId)
	if err != nil {
		log.Debugln("Failed conversion of layerId to pointer ", err)
		return err
	}

	infop, err := convertInfo(info)
	if err != nil {
		log.Debugln("Failed conversion of driver info ", err)
		return err
	}

	// Call the procedure itself.
	r1, _, _ := procUnprepareLayer.Call(
		uintptr(unsafe.Pointer(&infop)),
		uintptr(unsafe.Pointer(layerIdp)))
	use(unsafe.Pointer(&infop))
	use(unsafe.Pointer(layerIdp))

	if r1 != 0 {
		return syscall.Errno(r1)
	}

	return nil
}
