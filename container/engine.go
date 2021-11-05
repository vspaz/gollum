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

func setHostname() {
	err := syscall.Sethostname([]byte("vspazzz"))
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

func mountFs() {
	err := syscall.Chroot("/home/vspaz/ubuntufs")
	if err != nil {
		panic(err)
	}
}

func changeIntoDirectory() {
	err := syscall.Chdir("/")
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

func child() {
	logger.Infof("Running %v\n", os.Args[2:])
	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd = setStdInOut(cmd)
	setHostname()
	mountProc()
	mountFs()
	changeIntoDirectory()
	runCommand(cmd)
	unmountProc()
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
