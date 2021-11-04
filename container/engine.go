package container

import (
	"github.com/sirupsen/logrus"
	"log"
	"os"
)

func Run() {
	logger := ConfigureLogger(logrus.InfoLevel)
	Dispatch(os.Args, logger)
}

func parent(logger *logrus.Logger) {

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
