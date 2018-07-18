package app // import "github.com/uesyn/kubernetes-external-gpu-scheduler/cmd/app"
import (
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/uesyn/kubernetes-external-gpu-scheduler/pkg/prioritizer"
	"github.com/uesyn/kubernetes-external-gpu-scheduler/util/logs"
)

func AddPrioritizeRoute(router *httprouter.Router, targetresource string) {
	p := prioritizer.NewExtendedResourcePrioritizer(targetresource)
	ph := PrioritizeHandler(p)
	router.POST("/prioritize", ph)
}

func Serve(port int, targetresource string) {
	router := httprouter.New()
	AddPrioritizeRoute(router, targetresource)
	logs.Infoln("Server starting on the port:", port)
	portstr := ":" + strconv.Itoa(port)
	if err := http.ListenAndServe(portstr, router); err != nil {
		logs.Errorln(err)
	}
}
