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

func child(logger *logrus.Logger) {

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
