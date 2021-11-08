package container

import (
	"os"
	"os/exec"
	"syscall"
)

func fork() {
	// execute yourself, /proc/self/exe is the copy of the binary image of the caller itself
	cmd := exec.Command("/proc/self/exe", append([]string{"subprocess"}, os.Args[2:]...)...)
	container := Container{cmd: cmd}
	container.setStdStreams()
	container.setNamespaces()
	container.runCommand()
}

func NewContainer(args []string) *Container {
	return &Container{cmd: exec.Command(args[2], args[3:]...)}
}

type Container struct {
	cmd *exec.Cmd
}

func (c *Container) setStdStreams() {
	c.cmd.Stdin = os.Stdin
	c.cmd.Stdout = os.Stdout
	c.cmd.Stderr = os.Stderr
}

func (c Container) setHostname(hostname string) {
	err := syscall.Sethostname([]byte(hostname))
	if err != nil {
		logger.Errorf("failed to set hostname %s", err)
	}
}

func (c Container) mountProc() {
	err := syscall.Mount("proc", "proc", "proc", 0, "")
	if err != nil {
		logger.Errorf("failed to mount proc %s", err.Error())
	}
}

func (c Container) mountFs(fsDirPath string) {
	err := syscall.Chroot(fsDirPath)
	if err != nil {
		panic(err)
	}
}

func (c Container) changeIntoDirectory(dirName string) {
	err := syscall.Chdir(dirName)
	if err != nil {
		logger.Errorf("failed to chdir directory: %s", err.Error())
	}
}

func (c *Container) runCommand() {
	if err := c.cmd.Run(); err != nil {
		logger.Panicf("failed to run cmd %s", err)
	}
}

func (c Container) unmountProc() {
	err := syscall.Unmount("proc", 0)
	if err != nil {
		logger.Errorf("failed to mount 'proc' %s", err.Error())
	}
}

func (c *Container) setNamespaces() {
	c.cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWNS | syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID,
	}
}

func subprocess() {
	logger.Infof("Running %v\n", os.Args[2:])
	container := NewContainer(os.Args)
	container.setStdStreams()
	container.setHostname("vspazzz")
	container.mountFs("/home/vspaz/ubuntufs")  // e.g. https://cdimage.ubuntu.com/ubuntu-base/releases/20.04/release/
	container.changeIntoDirectory("/")
	container.mountProc()
	container.runCommand()
	container.unmountProc()
}
