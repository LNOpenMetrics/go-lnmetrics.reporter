package plugin

import (
	"fmt"
	"github.com/LNOpenMetrics/go-lnmetrics.reporter/pkg/graphql"
	"github.com/vincenzopalazzo/glightning/glightning"
	"os"
)

// DevServerUpload dev command, need to be runned only in the intergation testing
func (instance *MetricOne) DevServerUpload(client *graphql.Client, lightning *glightning.Lightning) error {
	// TODO move this a compiled variable, to establish this a runtime
	if os.Getenv("DEVELOPER") == "" {
		return fmt.Errorf("not in dev mode")
	}
	return instance.UploadOnRepo(client, lightning)
}
