package prioritizer

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/uesyn/kubernetes-external-gpu-scheduler/util/logs"

	schedulerapi "k8s.io/kubernetes/pkg/scheduler/api"
)

func checkBody(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}
}

func AddPrioritizeRoute(prioritizer Prioritizer) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		checkBody(w, r)

		var buf bytes.Buffer
		body := io.TeeReader(r.Body, &buf)
		log.Print("info: ", prioritizer.Name, " ExtenderArgs = ", buf.String())

		var extenderArgs schedulerapi.ExtenderArgs
		var hostPriorityList *schedulerapi.HostPriorityList

		if err := json.NewDecoder(body).Decode(&extenderArgs); err != nil {
			panic(err)
		}

		if list, err := prioritizer.Handler(extenderArgs); err != nil {
			panic(err)
		} else {
			hostPriorityList = list
		}

		if resultBody, err := json.Marshal(hostPriorityList); err != nil {
			panic(err)
		} else {
			log.Print("info: ", prioritizer.Name, " hostPriorityList = ", string(resultBody))
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(resultBody)
		}
	}
}

func LoggingServer(handle httprouter.Handle, path string) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		logs.Debugln("PATH: ", path, "REQUEST: ", r.Body)
		handle(w, r, p)
		logs.Debugln("PATH: ", path, "RESPONSE: ", w)
	}
}

func AddPrioritize(router *httprouter.Router, prioritizer Prioritizer) {
	path := prioritiesPrefix + "/" + prioritizer.Name
	router.POST(path, LoggingServer(AddPrioritizeRoute(prioritizer), path))
}
