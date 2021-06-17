package plugin

import (
	"errors"
	"fmt"
)

type MemoryDB struct {
	// Node id of the actual node
	nodeId string
	// List for metrics that the plugin need to collect
	metrics map[int]*Metrics
}

func (instance *MemoryDB) AddNodeId(id *string) {
	instance.nodeId = *id
}

func (instance *MemoryDB) RegisterMetrics(id int, metric *Metrics) error {
	_, ok := instance.metrics[id]
	if ok {
		return errors.New(fmt.Sprintf("A metrics with the following id %d is already taken", id))
	}
	instance.metrics[id] = metric
	return nil
}

func (instance *MemoryDB) GetMetrics(id int) (*Metrics, error) {
	value, ok := instance.metrics[id]
	if !ok {
		return nil, errors.New(fmt.Sprintf("A metrics with the following id %id it is not registered", id))
	}
	return value, nil
}

func (instance *MemoryDB) MakePersistent() error {
	for _, metric := range instance.metrics {
		err := (*metric).MakePersistent()
		if err != nil {
			return err
		}
	}
	return nil
}
