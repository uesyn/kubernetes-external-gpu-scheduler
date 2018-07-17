package app // import "github.com/uesyn/kubernetes-external-gpu-scheduler/cmd/app"
import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/uesyn/kubernetes-external-gpu-scheduler/pkg/prioritizer"
	"github.com/uesyn/kubernetes-external-gpu-scheduler/util/logs"
)

func Serve(port int, targetresource string) {
	router := httprouter.New()
	p := prioritizer.NewExtendedResourcePrioritizer(targetresource)
	ph := PrioritizerHandler(p)
	router.POST("prioritize", ph)
	logs.Info("Server starting on the port :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		logs.Errorln(err)
	}
}
