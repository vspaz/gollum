package container

import (
	"os"
)

func Dispatch(args []string) {
	switch args[1] {
	case "run":
		fork()
	case "subprocess":
		subprocess()
	default:
		logger.Panicf("method undefined %s", args[1])
	}
}

func Run() {
	Dispatch(os.Args)
}
