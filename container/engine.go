package container

import (
	"os"
)

// e.g. https://cdimage.ubuntu.com/ubuntu-base/releases/20.04/release/
const imagePath = "/home/vspaz/ubuntufs"

func Dispatch(args []string) {
	switch args[1] {
	case "run":
		fork()
	case "subprocess":
		subprocess(imagePath)
	default:
		logger.Panicf("method undefined %s", args[1])
	}
}

func Run() {
	Dispatch(os.Args)
}
