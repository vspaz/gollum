package container

import (
	"github.com/sirupsen/logrus"
	"log"
	"os"
	"os/exec"
	"syscall"
)

func Run() {
	logger := ConfigureLogger(logrus.InfoLevel)
	Dispatch(os.Args, logger)
}

func setStdInOut(cmd *exec.Cmd) *exec.Cmd {
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd
}

func parent(logger *logrus.Logger) {
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

func setHostname(logger *logrus.Logger) {
	err := syscall.Sethostname([]byte("vspazzz"))
	if err != nil {
		logger.Errorf("failed to set hostname %s", err)
	}
}

func child(logger *logrus.Logger) {
	logger.Infof("Running %v\n", os.Args[2:])
	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd = setStdInOut(cmd)
	setHostname(logger)

	syscall.Mount("proc", "proc", "proc", 0, "")
	err := syscall.Chroot("/home/vspaz/ubuntufs")
	if err != nil {
		panic(err)
	}
	syscall.Chdir("/")
	if err := cmd.Run(); err != nil {
		panic(err)
	}
	syscall.Unmount("proc", 0)
}

func Dispatch(args []string, logger *logrus.Logger) {
	switch args[1] {
	case "run":
		parent(logger)
	case "child":
		child(logger)
	default:
		log.Panicf("method undefined %s", args[1])
	}
}
