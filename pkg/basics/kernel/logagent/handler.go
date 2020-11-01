package logagent

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/vmihailenco/msgpack"
	// "strings"
	"net"
	// "net/http"
	"os"
	"sync"
)

var (
	// remoteAddr string = "10.25.96.167:9123"
	remoteAddr string = "127.0.0.1:9123"
)

var (
	AGENT_MAGIC       string = "emit"
	AGETN_MODLE_COUNT string = "counter"
	AGENT_MODLE_STORE string = "store"
	AGENT_MODLE_TIMER string = "timer"
)

var (
	STATUS_CODE string = "code"
	METHOD      string = "event"
	REMOTE_PS   string = "remote_ps"
	LOCAL_HOST  string = "local_host"
	REMOTE_HOST string = "remote_host"

	PROJECT string = "project"
	PORT    string = "port"

	MEM_CACHE   string = "mem_cache"
	REDIS_CACHE string = "redis_cache"

	INFCLIENT string = "clientTag"
)

var (
	NO_FOUND_AVALIABLE_CONN = errors.New("no use conn,check logagent server is ok")
)

type LogAgent struct {
	lock     *sync.RWMutex
	conn     *Conn
	endPoint string
}

type LogAgentTagOptionsInterface interface {
	SetRemoteAddr(remoteAddr string)
	GetRemoteAddr() string

	SetTag(key string, value string)
	GetTag(key string) string

	SetMemCache(value string)
	GetMemCache() string

	GetValues() map[string]string
}

type LogAgentTagOptions struct {
	valueMap map[string]string
}

func NewLogAgentTagOptions() *LogAgentTagOptions {
	return &LogAgentTagOptions{
		valueMap: make(map[string]string),
	}
}

func (lg *LogAgentTagOptions) SetMemCache(value string) {

}
func (lg *LogAgentTagOptions) GetMemCache() string {

	if value, ok := lg.valueMap[MEM_CACHE]; ok {
		return value
	}
	return ""
}

func (lg *LogAgentTagOptions) GetValues() map[string]string {
	return lg.valueMap
}

func (lg *LogAgentTagOptions) SetRemoteAddr(remoteAddr string) {
	lg.valueMap[REMOTE_HOST] = remoteAddr
}

func (lg *LogAgentTagOptions) GetRemoteAddr() string {

	if value, ok := lg.valueMap[REMOTE_HOST]; ok {
		return value
	}
	return ""
}

func (lg *LogAgentTagOptions) SetTag(key string, value string) {
	lg.valueMap[key] = value
}

func (lg *LogAgentTagOptions) GetTag(key string) string {
	if value, ok := lg.valueMap[key]; ok {
		return value
	}
	return ""
}

// GetLocalIP 获取本机IP
func getLocalIP() ([]string, error) {
	ret := []string{}
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ret, err
	}
	for _, a := range addrs {
		if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				ret = append(ret, ipnet.IP.String())
			}
		}
	}
	return ret, err
}

func InitAgent() *LogAgent {

	defer func() {
		recover()
	}()

	lock := new(sync.RWMutex)

	conn, err := NewLogConn(remoteAddr)

	if err != nil {
		conn = nil
	}

	localIP, _ := getLocalIP()
	endPoint, _ := os.Hostname()
	if len(endPoint) == 0 {
		endPoint = localIP[0]
	}

	agent := &LogAgent{
		lock:     lock,
		conn:     conn,
		endPoint: endPoint,
	}
	return agent
}

func (la *LogAgent) sendMsg(msg []byte) error {

	la.lock.Lock()
	defer la.lock.Unlock()

	conn := la.conn

	if conn == nil || conn.Closed() {
		newConn, err := NewLogConn(remoteAddr)
		if err != nil {
			return NO_FOUND_AVALIABLE_CONN
		}
		la.conn = newConn
		conn = newConn
	}
	conn.Write(msg, 0)
	return nil
}

func agentsubstr(str string, start int, end int) string {
	rs := []rune(str)
	length := len(rs)

	if start < 0 || start > length {
		panic("start is wrong")
	}

	if end < 0 || end > length {
		panic("end is wrong")
	}

	return string(rs[start:end])
}

func (la *LogAgent) makeValueInfo(value map[string]string) string {

	var rsp string

	// value[LOCAL_HOST] = la.endPoint

	for key, value := range value {

		rsp = rsp + key + "=" + value + "|"
	}

	if len(rsp) < 3 {
		return ""
	}

	v := agentsubstr(rsp, len(rsp)-1, len(rsp))
	rspValue := rsp
	if v == "|" {
		rspValue = agentsubstr(rsp, 0, len(rsp)-1)
	}
	return rspValue
}

func makeMspPackMsg(msgArr []string) ([]byte, error) {
	var buf bytes.Buffer
	enc := msgpack.NewEncoder(&buf)
	err := enc.Encode(msgArr)
	if err != nil {
		var tmp []byte
		return tmp, err
	}
	return buf.Bytes(), nil
}

func (la *LogAgent) makeTimerLog(service string, timecost int64, value map[string]string) ([]byte, error) {
	var msgArr []string
	// 【  emit,timer,service,timecost 】

	msgArr = append(msgArr, AGENT_MAGIC)
	msgArr = append(msgArr, AGENT_MODLE_TIMER)

	msgArr = append(msgArr, service)

	timecostStr := fmt.Sprintf("%d", timecost)
	msgArr = append(msgArr, timecostStr)

	statValue := la.makeValueInfo(value)
	msgArr = append(msgArr, statValue)

	msgArr = append(msgArr, "")

	// fmt.Println("makeTimerLog:", msgArr)

	return makeMspPackMsg(msgArr)
}

func (la *LogAgent) makeSnapshotLog(metric string, snapshotValue float64, value map[string]string) ([]byte, error) {

	var msgArr []string
	msgArr = append(msgArr, AGENT_MAGIC)
	msgArr = append(msgArr, AGENT_MODLE_STORE)
	msgArr = append(msgArr, metric)

	snapshotStr := fmt.Sprintf("%f", snapshotValue)
	msgArr = append(msgArr, snapshotStr)

	statValue := la.makeValueInfo(value)
	msgArr = append(msgArr, statValue)
	msgArr = append(msgArr, "")
	// fmt.Println("store msgArr:", msgArr)
	return makeMspPackMsg(msgArr)
}

func (la *LogAgent) getValueMap(lat LogAgentTagOptionsInterface) map[string]string {

	var valueMap map[string]string
	if lat == nil {
		lat = NewLogAgentTagOptions()
	}
	valueMap = lat.GetValues()

	return valueMap
}

func (la *LogAgent) SendSnapshotLog(metric string, snapshotValue float64, lat LogAgentTagOptionsInterface) error {

	defer func() {
		recover()
	}()

	valueMap := la.getValueMap(lat)

	return la.sendSnapshotLog(metric, snapshotValue, valueMap)
}

func (la *LogAgent) sendSnapshotLog(metric string, snapshotValue float64, valueMap map[string]string) error {

	msgByte, err := la.makeSnapshotLog(metric, snapshotValue, valueMap)
	if err != nil {
		return err
	}
	err = la.sendMsg(msgByte)
	return err
}

func (la *LogAgent) SendGoroutineCount(service string, count int) error {
	metric := "event.code.gocount"
	var snapshotValue float64
	snapshotValue = float64(count)

	valueMap := make(map[string]string)
	valueMap[LOCAL_HOST] = la.endPoint
	valueMap[PROJECT] = service

	return la.sendSnapshotLog(metric, snapshotValue, valueMap)
}

func (la *LogAgent) SendServiceAlive(service string, port int) error {

	metric := "event.code.alive"
	var snapshotValue float64
	snapshotValue = 1
	valueMap := make(map[string]string)
	valueMap[LOCAL_HOST] = la.endPoint
	valueMap[PROJECT] = service
	valueMap[PORT] = fmt.Sprintf("%d", port)
	return la.sendSnapshotLog(metric, snapshotValue, valueMap)

}

func (la *LogAgent) SendTimerLog(service, remoteService, metric string, timecost int64, code int, lat LogAgentTagOptionsInterface) error {
	defer func() {
		recover()
	}()

	valueMap := la.getValueMap(lat)
	return la.sendTimerLog(service, remoteService, metric, timecost, code, valueMap)
}

func (la *LogAgent) sendTimerLog(service string, remoteService string, metric string, timecost int64, code int, valueMap map[string]string) error {
	defer func() {
		recover()
	}()

	valueMap[STATUS_CODE] = fmt.Sprintf("%d", code)
	valueMap[METHOD] = metric
	valueMap[REMOTE_PS] = remoteService
	msgByte, err := la.makeTimerLog(service, timecost, valueMap)

	if err != nil {
		return err
	}
	err = la.sendMsg(msgByte)
	return err
}

///stat直接写入
func (la *LogAgent) SendTimerLogMsg(service string, remoteService string, metric string, timecost int64, code int) error {

	valueMap := make(map[string]string)

	return la.sendTimerLog(service, remoteService, metric, timecost, code, valueMap)
}
