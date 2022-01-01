// This package contains all th code that hide the complexity
// of the db itself.
package db

// Plugin database interface
type PluginDatabase interface {
	// Wrapper around the method to store data
	// in the key value db
	PutValue(key string, value *string) error

	// Wrapper around the method to get data
	// in the key value database
	GetValue(key string) (*string, error)

	// Wrapper around the method to get data
	// in the key value database
	DeleteValue(key string) error

	// wrapper to check if the database it is ready
	IsReady() bool

	// Migrate procedure from one metrics to another, in
	// input we get the metric name that we want migrate.
	Migrate(metrics []*string) error

	// Load last not commitment metric status, given by
	// the last update
	LoadLastMetricOne() (*string, error)

	// store a snapshot of in the local database in a specify
	// moment defined with a UNIX timestamp.
	//
	// This will hide the logic under the database.
	StoreMetricOneSnapshot(timestamp int64, payload *string) error

	// get the information that are stored in the with old key, this
	// help to very hard migration of the database where the more easy
	// thinks to do is to store the information inside a "old" key and
	// start to from a clean key the new format of the metrics.
	//
	// @key: is the old key where the implementation is stored
	// @erase: it is a flag that tell to the db to erase the information
	// forever after the query.
	GetOldData(key string, erase bool) (*string, bool)

	// Store Metrics it is a generic method
	// that take a Metrics interface and store it
	// in th database.
	CloseDatabase() error

	// Erase database and all the data are lost forever
	EraseDatabase() error

	// Get location of the DB
	GetDBPath() string
}
