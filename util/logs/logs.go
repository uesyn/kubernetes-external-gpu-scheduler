package logs // import "github.com/uesyn/kubernetes-external-gpu-scheduler/util/logs"

import (
	"log"
)

func Infof(format string, a ...interface{}) {
	format = "info: " + format
	var b []interface{}
	for _, val := range a {
		b = append(b, val)
	}
	log.Printf(format, b)
}

func Info(a ...interface{}) {
	var b []interface{} = []interface{}{"info: "}
	for _, val := range a {
		b = append(b, val)
	}
	log.Print(b...)
}

func Infoln(a ...interface{}) {
	var b []interface{} = []interface{}{"info: "}
	for _, val := range a {
		b = append(b, val)
	}
	log.Println(b...)
}

func Tracef(format string, a ...interface{}) {
	format = "trace: " + format
	var b []interface{}
	for _, val := range a {
		b = append(b, val)
	}
	log.Printf(format, b)
}

func Trace(a ...interface{}) {
	var b []interface{} = []interface{}{"trace: "}
	for _, val := range a {
		b = append(b, val)
	}
	log.Print(b...)
}

func Traceln(a ...interface{}) {
	var b []interface{} = []interface{}{"trace: "}
	for _, val := range a {
		b = append(b, val)
	}
	log.Println(b...)
}

func Debugf(format string, a ...interface{}) {
	format = "debug: " + format
	var b []interface{}
	for _, val := range a {
		b = append(b, val)
	}
	log.Printf(format, b)
}

func Debug(a ...interface{}) {
	var b []interface{} = []interface{}{"debug: "}
	for _, val := range a {
		b = append(b, val)
	}
	log.Print(b...)
}

func Debugln(a ...interface{}) {
	var b []interface{} = []interface{}{"debug: "}
	for _, val := range a {
		b = append(b, val)
	}
	log.Println(b...)
}

func Warnf(format string, a ...interface{}) {
	format = "warn: " + format
	var b []interface{}
	for _, val := range a {
		b = append(b, val)
	}
	log.Printf(format, b)
}

func Warn(a ...interface{}) {
	var b []interface{} = []interface{}{"warn: "}
	for _, val := range a {
		b = append(b, val)
	}
	log.Print(b...)
}

func Warnln(a ...interface{}) {
	var b []interface{} = []interface{}{"warn: "}
	for _, val := range a {
		b = append(b, val)
	}
	log.Println(b...)
}

func Errorf(format string, a ...interface{}) {
	format = "error: " + format
	var b []interface{}
	for _, val := range a {
		b = append(b, val)
	}
	log.Printf(format, b)
}

func Error(a ...interface{}) {
	var b []interface{} = []interface{}{"error: "}
	for _, val := range a {
		b = append(b, val)
	}
	log.Print(b...)
}

func Errorln(a ...interface{}) {
	var b []interface{} = []interface{}{"error: "}
	for _, val := range a {
		b = append(b, val)
	}
	log.Println(b...)
}

func Alertf(format string, a ...interface{}) {
	format = "alert: " + format
	var b []interface{}
	for _, val := range a {
		b = append(b, val)
	}
	log.Printf(format, b)
}

func Alert(a ...interface{}) {
	var b []interface{} = []interface{}{"alert: "}
	for _, val := range a {
		b = append(b, val)
	}
	log.Print(b...)
}

func Alertln(a ...interface{}) {
	var b []interface{} = []interface{}{"alert: "}
	for _, val := range a {
		b = append(b, val)
	}
	log.Println(b...)
}
