package db

import (
	"fmt"
	"strconv"
	"strings"

	db "github.com/LNOpenMetrics/lnmetrics.utils/db/leveldb"
	"github.com/LNOpenMetrics/lnmetrics.utils/log"
)

type LevelDB struct {
	// Global database version to know basically
	// how the key are stored.
	dbVersion int
	// Global Key in the database to get the data model version
	dbKeyDb string
	// dictionary with the a diction that contains a mapping
	// to get the correct version by the following info:
	// - metric_name
	// - dbVersion
	metricsDbKeys map[string]string
}

func NewLevelDB(path string) (PluginDatabase, error) {

	dbKey := "db_data_version"

	if err := db.GetInstance().InitDB(path); err != nil {
		return nil, err
	}

	dataVersion, err := db.GetInstance().GetValue(dbKey)
	if err != nil {
		log.GetInstance().Error(fmt.Sprintf("Checking DV data version return an error like: %s", err))
		dataVersion = "1"
		if err := db.GetInstance().PutValue(dbKey, dataVersion); err != nil {
			return nil, err
		}
	}

	log.GetInstance().Info(fmt.Sprintf("Db with data model version %s", dataVersion))

	dataVersionConv, err := strconv.Atoi(dataVersion)

	if err != nil {
		return nil, err
	}

	return &LevelDB{
		dbVersion: dataVersionConv,
		dbKeyDb:   dbKey,
		metricsDbKeys: map[string]string{
			"metric_one/1": "metric_one",
			"metric_one/2": "metric_one",
		},
	}, nil
}

func (instance *LevelDB) PutValue(key string, value *string) error {
	return db.GetInstance().PutValue(key, *value)
}

func (instance *LevelDB) GetValue(key string) (*string, error) {
	value, err := db.GetInstance().GetValue(key)
	return &value, err
}

func (instance *LevelDB) DeleteValue(key string) error {
	return db.GetInstance().DeleteValue(key)
}

func (instance *LevelDB) IsReady() bool {
	return db.GetInstance().Ready()
}

func (instance *LevelDB) StoreMetricOneSnapshot(timestamp int, payload *string) error {
	key := strings.Join([]string{"metric_one", fmt.Sprint(timestamp)}, "/")
	if err := instance.PutValue(key, payload); err != nil {
		return err
	}
	timestampStr := fmt.Sprint(timestamp)
	keyLastUpt := strings.Join([]string{"metric_one", "last"}, "/")
	if err := instance.PutValue(keyLastUpt, &timestampStr); err != nil {
		return err
	}
	return nil
}

func (instance *LevelDB) LoadLastMetricOne() (*string, error) {
	keyValue := strings.Join([]string{"metric_one", "last"}, "/")
	lastUpdate, err := instance.GetValue(keyValue)
	if err != nil {
		return nil, fmt.Errorf("Last metric it is not present in the db")
	}
	return lastUpdate, nil
}

// Take the version of the data and apply the procedure
// to migrate the database.
func (instance *LevelDB) Migrate(metrics []*string) error {
	for _, metric := range metrics {
		switch *metric {
		case "metric_one":
			if err := instance.migrateMetricOne(); err != nil {
				return err
			}
		default:
			return fmt.Errorf("Metric with key %s is not supported", *metric)
		}
	}

	return nil
}

// Private function that migrate (if needed) the metrics from one key to another key.
func (instance *LevelDB) migrateMetricOne() error {
	// 1. Check the global version of the db
	// 2. Load the metric
	// 3. See if the metrics can be migrated by version number
	// 3.1 Migrate the version with in the new version of the data-model.
	if instance.dbVersion == 1 {
		// migrate to version two.
		return instance.migrateMetricOneToVersionTwo()
	}
	return nil
}

func (instance *LevelDB) migrateMetricOneToVersionTwo() error {
	dictKey := strings.Join([]string{"metric_one", fmt.Sprint(instance.dbVersion)}, "/")
	metricKey := instance.metricsDbKeys[dictKey]

	// the full payload it is store in the single
	// instance
	metricJson, err := db.GetInstance().GetValue(metricKey)
	if err != nil {
		return err
	}

	instance.dbVersion = 2
	dictKey = strings.Join([]string{"metric_one", fmt.Sprint(instance.dbVersion)}, "/")

	oldMetricKey := metricKey
	metricKey = instance.metricsDbKeys[dictKey]
	oldKey := strings.Join([]string{metricKey, "old"}, "/")
	if err := db.GetInstance().PutValue(oldKey, metricJson); err != nil {
		return err
	}

	if err := db.GetInstance().DeleteValue(oldMetricKey); err != nil {
		return err
	}

	if err := db.GetInstance().PutValue(instance.dbKeyDb, fmt.Sprint(instance.dbVersion)); err != nil {
		return err
	}

	// store the new version of the data.
	return nil
}

// Close the database
func (instance *LevelDB) CloseDatabase() error {
	return db.GetInstance().CloseDatabase()
}

// Erase the Database and lost the data forever
func (instance *LevelDB) EraseDatabase() error {
	return db.GetInstance().EraseDatabase()
}
