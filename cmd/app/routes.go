package prioritizer

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/uesyn/kubernetes-external-gpu-scheduler/pkg/prioritizer"
	"github.com/uesyn/kubernetes-external-gpu-scheduler/util/logs"

	schedulerapi "k8s.io/kubernetes/pkg/scheduler/api"
)

func checkBody(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		logs.Errorln("Please send a request body")
		http.Error(w, "Please send a request body", 400)
		return
	}
}

func PrioritizeHandler(prioritizer prioritizer.Prioritizer) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		checkBody(w, r)

		var buf bytes.Buffer
		body := io.TeeReader(r.Body, &buf)

		var extenderArgs schedulerapi.ExtenderArgs
		var hostPriorityList *schedulerapi.HostPriorityList

		if err := json.NewDecoder(body).Decode(&extenderArgs); err != nil {
			logs.Error(err)
		}

		if list, err := prioritizer.Prioritize(*extenderArgs.Pod, extenderArgs.Nodes); err != nil {
			logs.Error(err)
		} else {
			hostPriorityList = list
		}

		if resultBody, err := json.Marshal(hostPriorityList); err != nil {
			logs.Error(err)
		} else {
			logs.Infoln(prioritizer.Name, " hostPriorityList = ", string(resultBody))
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(resultBody)
		}
	}
}
