package app

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
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

		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			logs.Errorln(err)
		}
		buf := bytes.NewBuffer(b)

		logs.Traceln("Receive Data:", string(b))

		var extenderArgs schedulerapi.ExtenderArgs
		var hostPriorityList *schedulerapi.HostPriorityList

		if err := json.NewDecoder(buf).Decode(&extenderArgs); err != nil {
			logs.Error(err)
		}
		logs.Traceln("Json Decode Data:", string(b))

		if list, err := prioritizer.Prioritize(&extenderArgs.Pod, extenderArgs.Nodes.Items); err != nil {
			logs.Error(err)
		} else {
			hostPriorityList = list
		}

		if resultBody, err := json.Marshal(hostPriorityList); err != nil {
			logs.Error(err)
		} else {
			logs.Debugln("HostPriorityList:", string(resultBody))
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(resultBody)
		}
	}
}
