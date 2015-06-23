package runconfig

import (
	"fmt"
	"strings"

	"github.com/docker/docker/opts"
)

func parseNetMode(netMode string) (NetworkMode, error) {
	parts := strings.Split(netMode, ":")
	switch mode := parts[0]; mode {
	case "default":
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
	return nil
}
