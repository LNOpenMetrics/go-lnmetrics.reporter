package persistence

import (
	"errors"
	"os"

	log "github.com/OpenLNMetrics/go-lnmetrics.reporter/pkg/log"
)

func PrepareHomeDirectory(lightningPath string) (*string, error) {
	fileInfo, err := os.Stat(lightningPath)
	if err != nil {
		return nil, err
	}

	switch mode := fileInfo.Mode(); {
	case mode.IsDir():
		log.GetInstance().Debug("Home dir " + lightningPath)
		path := lightningPath + "/metrics"
		_, err := os.Stat(path)
		if !os.IsNotExist(err) {
			return &path, nil
		}
		err = os.Mkdir(path, 0755)
		if err != nil {
			return nil, err
		}
		log.GetInstance().Info("Created directory at " + path)
		return &path, nil

	default:
		return nil, errors.New("This is a file")
	}

}
