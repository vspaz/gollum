package container

import (
	"net"
	"os"
	"os/exec"
	"syscall"
	"time"
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

func (c *Container) awaitNetworkUp(interval int) {
	start := time.Now()
	maxWait := time.Second * time.Duration(interval)
	for {
		interfaces, err := net.Interfaces()
		if err != nil {
			logger.Fatal("error: %s", err.Error())
		}
		if len(interfaces) > 1 {
			break
		}
		if time.Since(start) > maxWait {
			logger.Fatal("error: timeout")
		}
		time.Sleep(time.Second)
	}
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
		Cloneflags: syscall.CLONE_NEWNS |
			syscall.CLONE_NEWUTS |
			syscall.CLONE_NEWPID |
			syscall.CLONE_NEWIPC |
			syscall.CLONE_NEWNET |
			syscall.CLONE_NEWUSER,
		UidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      os.Getegid(),
				Size:        1,
			},
		},
		GidMappings: []syscall.SysProcIDMap{
			{
				ContainerID: 0,
				HostID:      os.Getegid(),
				Size:        1,
			},
		},
	}
}

func subprocess(imagePath string) {
	logger.Infof("Running %v\n", os.Args[2:])
	container := NewContainer(os.Args)
	container.setStdStreams()
	container.setHostname("gollum")
	container.mountFs(imagePath)
	container.changeIntoDirectory("/")
	container.mountProc()
	container.awaitNetworkUp(5)
	container.runCommand()
	container.unmountProc()
}
