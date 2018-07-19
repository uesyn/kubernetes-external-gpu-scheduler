package cmd // import "github.com/uesyn/kubernetes-external-gpu-scheduler/cmd"

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/uesyn/kubernetes-external-gpu-scheduler/cmd/app"
	"github.com/uesyn/kubernetes-external-gpu-scheduler/cmd/options"
	"github.com/uesyn/kubernetes-external-gpu-scheduler/k8sclient"
	"github.com/uesyn/kubernetes-external-gpu-scheduler/util/logs"
)

var rootCmd = &cobra.Command{
	Use: "kube-extendedresource-scheduler",
	Long: `This is extended scheduler, implemented only prioritizer for kubernetes.
This can choice a node which is high usage of extended resource.
For example, in case that you want to cram pods requiring GPU resource into a high GPU resource usage node, this scheduler is useful.`,
	Run: entrypoint,
}

func init() {
	options.SetValue(options.HELP, rootCmd.Flags().BoolP("help", "h", false, "Show this help"))
	options.SetValue(options.PORT, rootCmd.Flags().IntP("port", "p", 8008, "Listen port"))
	options.SetValue(options.TARGET, rootCmd.Flags().StringP("target", "t", "nvidia.com/gpu", "Target Extended Resource"))
	options.SetValue(options.KUBECONFIG, rootCmd.Flags().StringP("kubeconfig", "k", "", "kubeconfig file path. If run in kubernetes pod, you can omit this option."))
	options.SetValue(options.LOGLEVEL, rootCmd.Flags().StringP("loglevel", "l", "info", "Log Level: trace, debug, info, warn, error, alert"))
}

func entrypoint(cmd *cobra.Command, args []string) {
	logs.SetMinLogLevel(options.GetLoglevel())
	if options.GetHelp() {
		cmd.Help()
		os.Exit(0)

	}
	options.Show()

	var err error
	if options.GetKubeConfig() == "" {
		err = k8sclient.SetConfigInCluster()
	} else {
		err = k8sclient.SetConfigFromKubeconfig(options.GetKubeConfig())
	}
	if err != nil {
		logs.Error(err)
	}
	app.Serve(options.GetPort(), options.GetTarget())
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		app.Serve(options.GetPort(), options.GetTarget())
		logs.Errorln(err)
	}
}
