package options

import "github.com/uesyn/kubernetes-external-gpu-scheduler/util/logs"

type Options struct {
	field map[Arg]interface{}
}

var options Options = Options{}

type Arg string

const (
	TARGET     = "target"
	PORT       = "port"
	LOGLEVEL   = "loglevel"
	HELP       = "help"
	KUBECONFIG = "kubeconfig"
)

func Show() {
	for key, value := range options.field {
		switch value.(type) {
		case *int:
			v, _ := options.field[key].(*int)
			logs.Infof("Command Argument: %s=%d\n", key, *v)
		case *string:
			v, _ := options.field[key].(*string)
			logs.Infof("Command Argument: %s=%s\n", key, *v)
		case *bool:
			v, _ := options.field[key].(*bool)
			logs.Infof("Command Argument: %s=%v\n", key, *v)
		default:
			logs.Infof("Command Argument: %s=%v\n", key, value)
		}
	}
}

func SetValue(name Arg, value interface{}) {
	if options.field == nil {
		options.field = make(map[Arg]interface{})
	}

	options.field[name] = value
}

func GetPort() int {
	v, ok := options.field[PORT].(*int)
	if ok {
		return *v
	}
	logs.Errorln("Invalid value:", options.field[PORT])
	return 0
}

func GetTarget() string {
	v, ok := options.field[TARGET].(*string)
	if ok {
		return *v
	}
	logs.Errorln("Invalid value:", options.field[TARGET])
	return ""
}

func GetKubeConfig() string {
	v, ok := options.field[KUBECONFIG].(*string)
	if ok {
		return *v
	}
	logs.Errorln("Invalid value:", options.field[KUBECONFIG])
	return ""
}

func GetLoglevel() string {
	v, ok := options.field[LOGLEVEL].(*string)
	if ok {
		return *v
	}
	logs.Errorln("Invalid value:", options.field[LOGLEVEL])
	return ""
}

func GetHelp() bool {
	v, ok := options.field[HELP].(*bool)
	if ok {
		return *v
	}
	logs.Errorln("Invalid value:", options.field[HELP])
	return false
}
