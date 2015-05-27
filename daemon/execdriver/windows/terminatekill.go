// +build windows

package windows

import (
	"github.com/Sirupsen/logrus"
	"github.com/docker/docker/daemon/execdriver"
	"github.com/docker/docker/pkg/hcsshim"
)

func (d *driver) Terminate(p *execdriver.Command) error {
	logrus.Debugf("WindowsExec: Terminate() id=%s", p.ID)
	return kill(p.ID, p.ContainerPid)
}

func (d *driver) Kill(p *execdriver.Command, sig int) error {
	logrus.Debugf("WindowsExec: Kill() id=%s sig=%d", p.ID, sig)
	return kill(p.ID, p.ContainerPid)
}

func kill(id string, pid int) error {
	logrus.Debugln("kill() ", id, pid)
	var err error

	// Terminate Process
	if err = hcsshim.TerminateProcessInComputeSystem(id, uint32(pid)); err != nil {
		logrus.Warnf("Failed to terminate pid %d in %s", id, pid, err)
		// Ignore errors
		err = nil
	}

	// Shutdown the compute system
	if err = hcsshim.ShutdownComputeSystem(id); err != nil {
		logrus.Errorf("Failed to shutdowm %s - %s", id, err)
	}

	return err
}
