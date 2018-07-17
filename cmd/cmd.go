package cmd // import "github.com/uesyn/kubernetes-external-gpu-scheduler/cmd"
import (
	"os"

	"github.com/spf13/cobra"
	"github.com/uesyn/kubernetes-external-gpu-scheduler/cmd/options"
	"github.com/uesyn/kubernetes-external-gpu-scheduler/util/logs"
	"mod/k8s.io/kubernetes@v1.10.0/cmd/controller-manager/app"
)

var rootCmd = &cobra.Command{
	Use: "kube-extendedresource-scheduler",
	Long: `This is extended scheduler, implemented only prioritizer for kubernetes.
This can choice a node which is high usage of extended resource.
For example, in case that you want to cram pods requiring GPU resource into a high GPU resource usage node, this scheduler is useful.`,
	Run: run,
}

func run(cmd *cobra.Command, args []string) {
	logs.SetMinLogLevel(options.GetLoglevel())
	if options.GetHelp() {
		cmd.Help()
		os.Exit(0)
	}
}

func initRootCmd() {
	options.SetValue("help", *rootCmd.Flags().BoolP("help", "h", false, "Show this help"))
	options.SetValue("port", *rootCmd.Flags().IntP("port", "p", 8008, "Listen port"))
	options.SetValue("target", *rootCmd.Flags().StringP("target", "t", "nvidia.com/gpu", "Target Extended Resource"))
	options.SetValue("loglevel", *rootCmd.Flags().StringP("loglevel", "l", "info", "Log Level: trace, debug, info, warn, error, alert"))
}

func Execute() {
	initRootCmd()

	if err := rootCmd.Execute(); err != nil {
		app.Serve(options.GetPort(), options.GetTarget())
		logs.Errorln(err)
	}
}
