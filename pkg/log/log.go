package log

import (
	"github.com/sirupsen/logrus"
	"os"
	"sync"
)

type singleton struct {
	Log *logrus.Logger
	sync.RWMutex
}

var instance singleton

func init() {
	//FIXME: give the possibility to specify the path.
	instance.Log = logrus.New()
	dirname, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	instance.Log.SetLevel(logrus.DebugLevel)
	file, err := os.OpenFile(dirname+"/metrics.log",
		os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		instance.Log.Out = file
	} else {
		instance.Log.Info("Failed to log to file, using default stderr")
	}

}

func GetInstance() *singleton {
	return &instance
}

func (this *singleton) Debug(message interface{}) {
	instance.Log.Debug(message)
}

func (this *singleton) Info(message interface{}) {
	instance.Log.Info(message)
}

func (this *singleton) Error(message interface{}) {
	instance.Log.Error(message)
}
