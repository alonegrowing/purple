package util

import "time"

var (
	ch = make(chan *StatEntry, 100000)
)

func (st *StatEntry) End(category string, code int) {
	st.code = code
	st.end = time.Now()
	st.category = category
	ch <- st
	// sendTimeLog(st)
}

type StatEntry struct {
	end           time.Time
	start         time.Time
	event         string
	category      string
	code          int
	remoteservice string
	splitpoint    string
	tag           string
}

func NewServiceStatEntry(client string, event string) *StatEntry {

	service := event
	if len(event) != 0 {
		event = client + "." + event
	} else {
		event = client
	}

	st := &StatEntry{
		event:         event,
		start:         time.Now(),
		remoteservice: service,
		splitpoint:    ".",
		tag:           client,
	}

	return st
}
