// +build !windows

package runconfig

import (
	"fmt"
	"strings"

	"github.com/docker/docker/opts"
)

func parseNetMode(netMode string) (NetworkMode, error) {
	parts := strings.Split(netMode, ":")
	switch mode := parts[0]; mode {
	case "default", "bridge", "none", "host":
	case "container":
		if len(parts) < 2 || parts[1] == "" {
			return "", fmt.Errorf("invalid container format container:<name|id>")
		}
	default:
		return "", fmt.Errorf("invalid --net: %s", netMode)
	}
	return NetworkMode(netMode), nil
}

func validateNetMode(netMode NetworkMode,
	flHostname *string,
	flLinks opts.ListOpts,
	flDns opts.ListOpts,
	flExtraHosts opts.ListOpts,
	flMacAddress *string,
	flPublish opts.ListOpts,
	flPublishAll *bool,
	flExpose opts.ListOpts) error {

	if (netMode.IsHost() || netMode.IsContainer()) && *flHostname != "" {
		return ErrConflictNetworkHostname
	}

	if netMode.IsHost() && flLinks.Len() > 0 {
		return ErrConflictHostNetworkAndLinks
	}

	if netMode.IsContainer() && flLinks.Len() > 0 {
		return ErrConflictContainerNetworkAndLinks
	}

	if (netMode.IsHost() || netMode.IsContainer()) && flDns.Len() > 0 {
		return ErrConflictNetworkAndDns
	}

	if (netMode.IsContainer() || netMode.IsHost()) && flExtraHosts.Len() > 0 {
		return ErrConflictNetworkHosts
	}

	if (netMode.IsContainer() || netMode.IsHost()) && *flMacAddress != "" {
		return ErrConflictContainerNetworkAndMac
	}

	if netMode.IsContainer() && (flPublish.Len() > 0 || *flPublishAll == true) {
		return ErrConflictNetworkPublishPorts
	}

	if netMode.IsContainer() && flExpose.Len() > 0 {
		return ErrConflictNetworkExposePorts
	}
	return nil
}
