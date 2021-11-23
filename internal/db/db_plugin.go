// This package contains all th code that hide the complexity
// of the db itself.
package db

import (
	"github.com/LNOpenMetrics/go-lnmetrics.reporter/internal/plugin"
)

// Plugin database interface
type PluginDatabase interface {
	// Create metric it is the interface that
	// Prepare the database for working with the plugin
	CreateMetric() error
	// Migrate procedure, that give the possibility to
	// migrate from one data-model to another.
	Migrate() error
	// GetDataVersion give the current value of the data-model
	// version in the db.
	GetDataVersion() (uint, error)
	// MetricBaseID give the base id of a metric, that in
	// noSQL database can be a free key, or nothing in a
	// Relational Database, because it managed the id itself.
	// What about the graph database?
	MetricBaseID(metric plugin.Metric) string
	// Store Metrics it is a generic method
	// that take a Metrics interface and store it
	// in th database.
	StoreMetric(metric plugin.Metric) error
	// GetMetric return a metric that it is stored in the db
	// for a period range [start, end], to disable the range
	// it is possible use the -1 value
	GetMetric(metricID uint, startPeriod int, endPeriod int) (plugin.Metric, error)
}
