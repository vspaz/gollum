package container

import (
	"log"
	"os"
)

func Run() {
	Dispatch(os.Args)
}

func parent() {

}

func child() {

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
