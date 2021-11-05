package container

import (
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"os/exec"
	"syscall"
)

var logger = ConfigureLogger(logrus.InfoLevel)

func Run() {
	Dispatch(os.Args)
}

func setStdInOut(cmd *exec.Cmd) *exec.Cmd {
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd
}

func parent() {
	logger.Infof("Running %v\n", os.Args[2:])
	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)
	cmd = setStdInOut(cmd)
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID,
	}

	if err := cmd.Run(); err != nil {
		logger.Panicf("error: %s", err.Error())
	}
}

func setHostname(hostname string) {
	err := syscall.Sethostname([]byte(hostname))
	if err != nil {
		logger.Errorf("failed to set hostname %s", err)
	}
}

func mountProc() {
	err := syscall.Mount("proc", "proc", "proc", 0, "")
	if err != nil {
		logger.Errorf("failed to mount proc %s", err.Error())
	}
}

func mountFs(fsDirPath string) {
	err := syscall.Chroot(fsDirPath)
	if err != nil {
		panic(err)
	}
}

func changeIntoDirectory(dirName string) {
	err := syscall.Chdir(dirName)
	if err != nil {
		logger.Errorf("failed to chdir directory: %s", err.Error())
	}
}

func unmountProc() {
	err := syscall.Unmount("proc", 0)
	if err != nil {
		logger.Errorf("failed to mount 'proc' %s", err.Error())
	}
}

func runCommand(cmd *exec.Cmd) {
	if err := cmd.Run(); err != nil {
		logger.Panicf("failed to run cmd %s", err)
	}
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

func child() {
	logger.Infof("Running %v\n", os.Args[2:])
	container := NewContainer(os.Args)
	container.setStdStreams()
	container.setHostname("vspazzz")
	container.mountProc()
	container.mountFs("/home/vspaz/ubuntufs")
	container.changeIntoDirectory("/")
	container.runCommand()
	container.unmountProc()
}

func Dispatch(args []string) {
	switch args[1] {
	case "run":
		parent()
	case "child":
		child()
	default:
		log.Panicf("method undefined %s", args[1])
	}
}
