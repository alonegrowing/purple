package sql

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"sync/atomic"
)

type Client struct {
	*gorm.DB
}

type Group struct {
	name    string
	master  *Client
	replica []*Client
	next    uint64
	total   uint64
}

func NewGroup(name string, master string, slaves []string) (*Group, error) {
	db, err := gorm.Open("mysql", master)
	if err != nil {
		return nil, fmt.Errorf("open mysql [%s] master %s error %s", name, master, err)
	}
	db = db.Debug()

	g := Group{
		name:    name,
		master:  &Client{DB: db},
		replica: make([]*Client, 0, len(slaves)),
		total:   0,
	}
	for _, slave := range slaves {
		db, err := gorm.Open("mysql", slave)
		if err != nil {
			return nil, fmt.Errorf("open mysql [%s] slave at %s error %s", name, slave, err)
		}
		db = db.Debug()
		g.replica = append(g.replica, &Client{
			DB: db,
		})
		g.total++
	}
	return &g, nil
}

func (g *Group) Master() *Client {
	return g.master
}

func (g *Group) Slave() *Client {
	if g.total == 0 {
		return g.master
	}
	next := atomic.AddUint64(&g.next, 1)
	return g.replica[next%g.total]
}

func (g *Group) Instance(isMaster bool) *Client {
	if isMaster {
		return g.Master()
	}
	return g.Slave()
}
