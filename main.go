package main // import "github.com/uesyn/kubernetes-external-gpu-scheduler"

import (
	"log"

	"github.com/comail/colog"
)

func main() {
	colog.SetDefaultLevel(colog.LDebug)
	colog.SetMinLevel(colog.LTrace)
	colog.SetFormatter(&colog.StdFormatter{
		Colors: true,
		Flag:   log.Ldate | log.Ltime | log.Lshortfile,
	})
	colog.Register()

	//	cmd.Execute()
}
