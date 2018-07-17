package options

import "github.com/uesyn/kubernetes-external-gpu-scheduler/util/logs"

type Options struct {
	field    map[string]interface{}
	port     int
	resource string
	loglevel string
	help     bool
}

var (
	options = &Options{}
)

func SetValue(name string, value interface{}) {
	options.field[name] = value
}

func GetPort() int {
	v, ok := options.field["port"].(int)
	if ok {
		return v
	}
	logs.Errorln("Invalid value:", options.field["port"])
	return 0
}

func GetTarget() string {
	v, ok := options.field["target"].(string)
	if ok {
		return v
	}
	logs.Errorln("Invalid value:", options.field["target"])
	return ""
}

func GetLoglevel() string {
	v, ok := options.field["loglevel"].(string)
	if ok {
		return v
	}
	logs.Errorln("Invalid value:", options.field["loglevel"])
	return ""
}

func GetHelp() bool {
	v, ok := options.field["help"].(bool)
	if ok {
		return v
	}
	logs.Errorln("Invalid value:", options.field["help"])
	return false
}
