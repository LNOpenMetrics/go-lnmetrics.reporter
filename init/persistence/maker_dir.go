package persistence

import (
	"errors"
	"os"

	log "github.com/OpenLNMetrics/go-metrics-reported/pkg/log"
)

func PrepareHomeDirectory(lightningPath string) error {
	fileInfo, err := os.Stat(lightningPath)
	if err != nil {
		return err
	}

	switch mode := fileInfo.Mode(); {
	case mode.IsDir():
		log.GetInstance().Debug("Home dir " + lightningPath)
		// TODO: Make the dir and if it is exist make it safe
	case mode.IsRegular():
		return errors.New("This is a file")
	}

	return nil

}
