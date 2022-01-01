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

	// Database path
	path string
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

	log.GetInstance().Info(fmt.Sprintf("DB data version: %d", dataVersionConv))
	return &LevelDB{
		dbVersion: dataVersionConv,
		dbKeyDb:   dbKey,
		path:      strings.Join([]string{path, "db"}, "/"),
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

func (instance *LevelDB) GetDBPath() string {
	return instance.path
}

func (instance *LevelDB) StoreMetricOneSnapshot(timestamp int64, payload *string) error {
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

	lastSnapshot := strings.Join([]string{"metric_one", *lastUpdate}, "/")
	metricJson, err := instance.GetValue(lastSnapshot)
	if err != nil {
		return nil, err
	}
	return metricJson, nil
}

func (instance *LevelDB) GetOldData(key string, erase bool) (*string, bool) {
	// Get the key of the prev version
	dictKey := strings.Join([]string{key, fmt.Sprint(instance.dbVersion - 1)}, "/")
	log.GetInstance().Info(fmt.Sprintf("Old data key: %s", dictKey))
	metricKey, found := instance.metricsDbKeys[dictKey]

	if !found {
		log.GetInstance().Info(fmt.Sprintf("No old key found in the mapping: key=%s/{db version -1}", key))
		return nil, false
	}

	oldKey := strings.Join([]string{metricKey, "old"}, "/")
	log.GetInstance().Info(fmt.Sprintf("Retrieval old metric with key: %s", oldKey))
	// the full payload it is store in the single
	// instance
	metricJson, err := db.GetInstance().GetValue(oldKey)
	if err != nil {
		log.GetInstance().Error(fmt.Sprintf("Error: %s", err))
		return nil, false
	}

	if erase {
		log.GetInstance().Infof("Erase old data on db with key: %s", oldKey)
		if err := db.GetInstance().DeleteValue(oldKey); err != nil {
			log.GetInstance().Error(fmt.Sprintf("Error: %s", err))
			return nil, false
		}
	}

	return &metricJson, true
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
		log.GetInstance().Info("No migration performed, because there is no data to migrate")
		return nil
	}

	instance.dbVersion = 2
	dictKey = strings.Join([]string{"metric_one", fmt.Sprint(instance.dbVersion)}, "/")

	oldMetricKey := metricKey
	metricKey = instance.metricsDbKeys[dictKey]
	oldKey := strings.Join([]string{metricKey, "old"}, "/")
	log.GetInstance().Debug(fmt.Sprintf("Storing old metric with key: %s", oldKey))
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
