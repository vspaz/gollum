package container

import (
	"log"
	"os"
)


func Run() {
	Dispatch(os.Args)
}

func Dispatch(args []string) {
	switch args[1] {
	case "run":
		fork()
	case "subprocess":
		subprocess()
	default:
		log.Panicf("method undefined %s", args[1])
	}
}
