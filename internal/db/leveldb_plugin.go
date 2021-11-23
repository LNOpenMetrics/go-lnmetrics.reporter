package db

import (
	_ "github.com/LNOpenMetrics/lnmetrics.utils/db/leveldb"

	"github.com/LNOpenMetrics/go-lnmetrics.reporter/internal/plugin"
)

type LevelDB struct {
	// The version of the data-model on the database
	VersionData uint
}

func NewLevelDB(path string) (PluginDatabase, error) {
	return &LevelDB{
		VersionData: 1,
	}, nil
}

// Empty implementation, here level db have no real data-model
func (instance *LevelDB) CreateMetric() error {
	// TODO: make sure that the index db will be there
	// TODO: make sure that the data version will be there
	return nil
}

// Take the version of the data and apply the procedure
// to migrate the database.
func (instance *LevelDB) Migrate() error {
	return nil
}

// Return the version of data-model
func (instance *LevelDB) GetDataVersion() (uint, error) {
	return instance.VersionData, nil
}

// TODO: Adding in the metric interface a new method that get back
// a map with all the information useful to make a metric id for database.
// P.S: We need a map because if we want migrate to a graphdb in this case
// we need another way to do that.
func (instance *LevelDB) MetricBaseID(metric plugin.Metric) string {
	return ""
}

// Store metric in the database,
// TODO: We need to specify how this process looks like,
// because we need to access to all the metrics details.
func (instance *LevelDB) StoreMetric(metric plugin.Metric) error {
	return nil
}

// Get the metric back
func (instance *LevelDB) GetMetric(metricID uint, startPeriod int, endPeriod int) (plugin.Metric, error) {
	// TODO: define how this method looks like, and how make this operation
	// with a generic metric.
	return nil, nil
}
