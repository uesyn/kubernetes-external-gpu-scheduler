package cmd // import "github.com/uesyn/kubernetes-external-gpu-scheduler/cmd"

import (
	"github.com/golang/glog"
	"github.com/spf13/cobra"
)

//var (
//	ZeroPriority = Prioritize{
//		Name: "zero_score",
//		Func: func(_ v1.Pod, nodes []v1.Node) (*schedulerapi.HostPriorityList, error) {
//			var priorityList schedulerapi.HostPriorityList
//			priorityList = make([]schedulerapi.HostPriority, len(nodes))
//			for i, node := range nodes {
//				priorityList[i] = schedulerapi.HostPriority{
//					Host:  node.Name,
//					Score: 0,
//				}
//			}
//			return &priorityList, nil
//		},
//	}
//)

//func main() {
//	level := StringToLevel(os.Getenv("LOG_LEVEL"))
//	log.Print("Log level was set to ", strings.ToUpper(level.String()))
//	colog.SetMinLevel(level)
//
//	router := httprouter.New()
//
//	priorities := []Prioritize{ZeroPriority}
//	for _, p := range priorities {
//		AddPrioritize(router, p)
//	}
//
//	log.Print("info: server starting on the port :8080")
//	if err := http.ListenAndServe(":8080", router); err != nil {
//		log.Fatal(err)
//	}
//}

type Options struct {
	port     int
	resource string
	loglevel string
}

var (
	options = &Options{}
)

var rootCmd = &cobra.Command{
	Use:   "kube-extendedresource-scheduler",
	Short: "Short comment",
	Long:  "Long comment",
}

func initRootCmd() {
	rootCmd.Flags().IntVarP(&options.port, "port", "p", 8008, "Listen port")
	rootCmd.Flags().StringVarP(&options.resource, "target", "t", "nvidia.com/gpu", "Target Extended Resource")
	rootCmd.Flags().StringVarP(&options.loglevel, "loglevel", "l", "info", "Log Level: trace, debug, info, warn, error, alert")
}

func Execute() {
	initRootCmd()

	if err := rootCmd.Execute(); err != nil {
		glog.Fatal(err)
	}
}
