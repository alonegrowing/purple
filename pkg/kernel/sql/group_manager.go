package sql

import (
	"fmt"
	"sync"
)

type GroupManager struct {
	mu     sync.RWMutex
	groups map[string]*Group
}

var (
	SQLGroupManager = newGroupManager()
)

func newGroupManager() *GroupManager {
	return &GroupManager{
		groups: make(map[string]*Group),
	}
}

func (gm *GroupManager) Add(name string, g *Group) error {
	gm.mu.Lock()
	defer gm.mu.Unlock()
	if _, ok := gm.groups[name]; ok {
		return fmt.Errorf("sql group alread exists")
	}
	gm.groups[name] = g
	return nil
}

func (gm *GroupManager) Get(name string) *Group {
	gm.mu.RLock()
	defer gm.mu.RUnlock()
	return gm.groups[name]
}

func (gm *GroupManager) PartitionBy(partiton func() (bool, string, string)) *Client {
	isMaster, dbName, tableName := partiton()
	return &Client{gm.Get(dbName).Instance(isMaster).Table(tableName)}
}

func Get(name string) *Group {
	return SQLGroupManager.Get(name)
}

func PartitionBy(partiton func() (bool, string, string)) *Client {
	return SQLGroupManager.PartitionBy(partiton)
}
