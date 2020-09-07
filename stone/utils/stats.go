package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"purple/stone/logagent"
	log "purple/stone/logging"
	"github.com/olekukonko/tablewriter"
)

var (
	ch = make(chan *StatEntry, 100000)
)

var (
	statUploadURL = "http://127.0.0.1:1988/v1/push"
	//statUploadURL = "http://10.45.9.153:1988/v1/push"
)

var eventTagMapping map[string]string
var eventPost map[string][]int // 上报的所有的数据

var (
	statMetric, _ = os.LookupEnv("STAT_METRIC")
	statPath, _   = os.LookupEnv("STAT_PATH")
)

var lg *logagent.LogAgent
var mLogger *log.Logger

type statInfo struct {
	max   int
	min   int
	avg   int
	count int
	total int
}
type tableDataSlice [][]string

func (t tableDataSlice) Len() int {
	return len(t)
}

func (t tableDataSlice) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
}

func (t tableDataSlice) Less(i, j int) bool {
	l := len(t[i])
	for m := 0; m < l; m++ {
		n := strings.Compare(t[i][m], t[j][m])
		if n < 0 {
			return true
		} else if n > 0 {
			return false
		}
	}
	return false
}

type statPostData struct {
	Metric      string  `json:"metric"`
	Endpoint    string  `json:"endpoint"`
	Timestamp   int64   `json:"timestamp"`
	Value       float64 `json:"value"`
	Step        int     `json:"step"`
	ContentType string  `json:"counterType"`
	Tags        string  `json:"tags"`
}

func SetStat(stFileName, stMetric string) {
	statPath = stFileName
	if stMetric != "" {
		statMetric = stMetric
	}
	mLogger.SetOutputByName(stFileName)
}

func init() {

	eventTagMapping = make(map[string]string)
	eventPost = make(map[string][]int)

	localIP, _ := GetLocalIP()
	endPoint, _ := os.Hostname()
	if len(endPoint) == 0 {
		endPoint = localIP[0]
	}
	//statPath, _ := os.LookupEnv("STAT_PATH")
	//statMetric, _ := os.LookupEnv("STAT_METRIC")
	statIntervalStr, _ := os.LookupEnv("STAT_INTERVAL")
	if statPath == "" {
		//statPath = "./stat.log"
	}
	if statMetric == "" {
		binaryName := strings.Split(os.Args[0], "/")
		statMetric = binaryName[len(binaryName)-1]
	}
	statInterval, _ := strconv.Atoi(statIntervalStr)
	if statInterval == 0 {
		statInterval = 60
	}
	mLogger = log.New()
	mLogger.SetFlags(0)
	mLogger.SetHighlighting(false)
	mLogger.SetRotateByDay()
	mLogger.SetPrintLevel(false)
	if statPath != "" {
		mLogger.SetOutputByName(statPath)
	}

	lg = logagent.InitAgent()

	statIntervalDuration := time.Duration(statInterval) * time.Second
	ticker := time.Tick(statIntervalDuration)
	stInfo := map[string]map[int]*statInfo{}
	runPid := os.Getpid()
	go func() {

		for {
			select {
			case <-ticker:
				statDts := []statPostData{}
				buf := bytes.NewBuffer(nil)
				now := time.Now()
				date := now.Format("2006-01-02 15:04:05")
				title := fmt.Sprintf("Server(Metric %s, Pid %d, LocalIP %s, Date %s) last %d seconds Statistic Info\n", statMetric, runPid, endPoint, date, statInterval)
				buf.WriteString(title)
				table := tablewriter.NewWriter(buf)
				table.SetHeader([]string{"Date", "Event", "Code", "Min(ms)", "Max(ms)", "Count", "Avg(ms)", "Rate(%)"})
				tableData := [][]string{}

				statEvent := make(map[string][]int)
				for k, v := range stInfo {
					k, v := k, v
					allTotal := 0
					for _, info := range v {
						allTotal += info.count
					}
					sdTotal := statPostData{
						Metric:      "event.total",
						Endpoint:    endPoint,
						Timestamp:   now.Unix(),
						Value:       float64(allTotal),
						Step:        statInterval,
						ContentType: "GAUGE",
						Tags:        fmt.Sprintf("project=%s,event=%s", statMetric, k),
					}
					statDts = append(statDts, sdTotal)
					for c, info := range v {
						//fmt.Printf("event %q category %q min %d max %d total %d count %d\n", k, c, info.min, info.max, info.total, info.count)
						rate := float64(info.count) / float64(allTotal) * 100
						avg := float64(info.total) / float64(info.count*1e6)

						var tags string

						if tag, ok := eventTagMapping[k]; ok {
							tags = fmt.Sprintf("project=%s,event=%s,code=%d,clienttag=%s", statMetric, k, c, tag)
						} else {
							tags = fmt.Sprintf("project=%s,event=%s,code=%d", statMetric, k, c)
						}

						if v, ok := eventPost[k]; ok {

							flag := checkInList(c, v)
							if flag == false {
								v = append(v, c)
								eventPost[k] = v
								// log.Debug("eventPost add, code_2:", c, ",k:", k)
							} else {
								// log.Debug("eventPost add, code_1:", c, ",k:", k)
							}
						} else {
							var codeList []int
							codeList = append(codeList, c)
							eventPost[k] = codeList
							// log.Debug("eventPost add, code:", c, ",k:", k)
						}

						if v, ok := statEvent[k]; ok {
							flag := checkInList(c, v)
							if flag == false {
								v = append(v, c)
								statEvent[k] = v
								// log.Debug("statEvent add, code_2:", c, ",k:", k)
							} else {
								// log.Debug("statEvent add, code_1:", c, ",k:", k)
							}

						} else {
							var codeList []int
							codeList = append(codeList, c)
							statEvent[k] = codeList
							// log.Debug("statEvent add, code:", c, ",k:", k)
						}

						sdRate := statPostData{
							Metric:      "event.code.rate",
							Endpoint:    endPoint,
							Timestamp:   now.Unix(),
							Value:       rate,
							Step:        statInterval,
							ContentType: "GAUGE",
							// Tags:        fmt.Sprintf("project=%s,event=%s,code=%d", statMetric, k, c),
							Tags: tags,
						}
						sdAvg := statPostData{
							Metric:      "event.code.avgtime",
							Endpoint:    endPoint,
							Timestamp:   now.Unix(),
							Value:       avg,
							Step:        statInterval,
							ContentType: "GAUGE",
							// Tags:        fmt.Sprintf("project=%s,event=%s,code=%d", statMetric, k, c),
							Tags: tags,
						}
						sdTotal := statPostData{
							Metric:      "event.code.count",
							Endpoint:    endPoint,
							Timestamp:   now.Unix(),
							Value:       float64(info.count),
							Step:        statInterval,
							ContentType: "GAUGE",
							// Tags:        fmt.Sprintf("project=%s,event=%s,code=%d", statMetric, k, c),
							Tags: tags,
						}
						statDts = append(statDts, sdRate, sdAvg, sdTotal)

						code := strconv.Itoa(c)
						min := fmt.Sprintf("%.4F", float64(info.min)/float64(1e6))
						max := fmt.Sprintf("%.4F", float64(info.max)/float64(1e6))
						count := strconv.Itoa(info.count)
						succ := fmt.Sprintf("%.2F", rate)
						tableData = append(tableData, []string{date, k, code, min, max, count, fmt.Sprintf("%.4F", avg), succ})
						//table.Append([]string{date, k, code, min, max, count, fmt.Sprintf("%.4F", avg), succ})
					}
				}

				// 判定本次没有上报的数据

				// log.Debug("eventPost:", eventPost, ",statEvent:", statEvent)

				for k, v := range eventPost {

					if statCode, ok := statEvent[k]; ok {

						for _, code := range v {

							flag := checkInList(code, statCode)

							if flag == false {

								// log.Debug("2 code ", code, " not in statCode:", statCode, " need send default,event:", k)

								event := k
								statDts = append(statDts, getDefaultAvgtime(event, endPoint, statInterval, code), getDefaultCount(event, endPoint, statInterval, code), getDefaultRate(event, endPoint, statInterval, code), getDefaultTotal(event, endPoint, statInterval, code))
							}
						}
						continue
					}
					for _, code := range v {

						// log.Debug("1 code ", code, " need send default,event:", k)

						event := k
						statDts = append(statDts, getDefaultAvgtime(event, endPoint, statInterval, code), getDefaultCount(event, endPoint, statInterval, code), getDefaultRate(event, endPoint, statInterval, code), getDefaultTotal(event, endPoint, statInterval, code))
						continue
					}
				}

				eventTagMapping = make(map[string]string)

				if len(tableData) == 0 {
					continue
				}
				sort.Sort(tableDataSlice(tableData))
				table.AppendBulk(tableData)
				table.Render()
				mLogger.Infof("%s", buf.String())
				stInfo = map[string]map[int]*statInfo{}
				httpClient := http.Client{
					Timeout: 3 * time.Second,
				}
				byteArrary, _ := json.Marshal(statDts)
				//fmt.Printf("post data %s\n", byteArrary)
				go func() {
					rsp, err := httpClient.Post(statUploadURL, "application/json", bytes.NewReader(byteArrary))
					if err != nil {
						log.GenLogf("upload stat info error %v, data %s", err, byteArrary)
						return
					}
					defer rsp.Body.Close()
					//io.Copy(ioutil.Discard, rsp.Body)
					data, _ := ioutil.ReadAll(rsp.Body)
					//fmt.Printf("post rsp %s\n", data)
					log.GenLogf("stat post %s, response %s", byteArrary, data)
				}()
			case st := <-ch:
				duration := int(st.end.Sub(st.start).Nanoseconds())

				key := getEventKey(st)
				if len(st.tag) != 0 {
					eventTagMapping[key] = st.tag
				}

				if _, ok := stInfo[key]; !ok {
					stInfo[key] = map[int]*statInfo{}
				}
				old, ok := stInfo[key][st.code]
				if !ok {
					stInfo[key][st.code] = &statInfo{max: duration, min: duration, total: duration, count: 1}
				} else {
					old.count++
					if duration > old.max {
						old.max = duration
					} else if duration < old.min {
						old.min = duration
					}
					old.total += duration
				}
			}
		}
	}()
}

func checkInList(c int, codeList []int) bool {

	var isExists bool
	isExists = false

	for _, code := range codeList {

		if code == c {
			isExists = true
			break
		}
	}
	return isExists
}

func getDefaultTotal(event, endpoint string, statInterval, code int) statPostData {

	allTotal := float64(0) / float64(1)

	now := time.Now()
	sdTotal := statPostData{
		Metric:      "event.total",
		Endpoint:    endpoint,
		Timestamp:   now.Unix(),
		Value:       float64(allTotal),
		Step:        statInterval,
		ContentType: "GAUGE",
		Tags:        fmt.Sprintf("project=%s,event=%s", statMetric, event),
	}
	return sdTotal
}

func getDefaultAvgtime(event, endpoint string, statInterval, code int) statPostData {
	now := time.Now()
	c := code

	var tags string
	if tag, ok := eventTagMapping[event]; ok {
		tags = fmt.Sprintf("project=%s,event=%s,code=%d,clienttag=%s", statMetric, event, c, tag)
	} else {
		tags = fmt.Sprintf("project=%s,event=%s,code=%d", statMetric, event, c)
	}

	avg := float64(0) / float64(1)

	sdAvg := statPostData{
		Metric:      "event.code.avgtime",
		Endpoint:    endpoint,
		Timestamp:   now.Unix(),
		Value:       avg,
		Step:        statInterval,
		ContentType: "GAUGE",
		Tags:        tags,
	}
	return sdAvg
}

func getDefaultRate(event, endpoint string, statInterval, code int) statPostData {
	now := time.Now()
	c := code
	var tags string
	if tag, ok := eventTagMapping[event]; ok {
		tags = fmt.Sprintf("project=%s,event=%s,code=%d,clienttag=%s", statMetric, event, c, tag)
	} else {
		tags = fmt.Sprintf("project=%s,event=%s,code=%d", statMetric, event, c)
	}

	rate := float64(0) / float64(1)

	sdRate := statPostData{
		Metric:      "event.code.rate",
		Endpoint:    endpoint,
		Timestamp:   now.Unix(),
		Value:       rate,
		Step:        statInterval,
		ContentType: "GAUGE",
		Tags:        tags,
	}
	return sdRate
}

func getDefaultCount(event, endpoint string, statInterval, code int) statPostData {
	now := time.Now()
	// date := now.Format("2006-01-02 15:04:05")
	var tags string
	c := code
	if tag, ok := eventTagMapping[event]; ok {
		tags = fmt.Sprintf("project=%s,event=%s,code=%d,clienttag=%s", statMetric, event, c, tag)
	} else {
		tags = fmt.Sprintf("project=%s,event=%s,code=%d", statMetric, event, c)
	}

	sdTotal := statPostData{
		Metric:      "event.code.count",
		Endpoint:    endpoint,
		Timestamp:   now.Unix(),
		Value:       float64(0),
		Step:        statInterval,
		ContentType: "GAUGE",
		Tags:        tags,
	}
	return sdTotal
}

func getEventKey(st *StatEntry) string {
	key := fmt.Sprintf("%s%s%s", st.event, st.splitpoint, st.category)
	return key
}

// GetLocalIP 获取本机IP
func GetLocalIP() ([]string, error) {
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

func (st *StatEntry) End(category string, code int) {
	st.code = code
	st.end = time.Now()
	st.category = category
	ch <- st
	// sendTimeLog(st)
}

func sendTimeLog(st *StatEntry) {

	key := getEventKey(st)
	duration := st.end.Sub(st.start).Nanoseconds() / 1e6

	tagOptional := logagent.NewLogAgentTagOptions()
	if len(st.tag) != 0 {
		tagOptional.SetTag(logagent.INFCLIENT, st.tag)
	}
	lg.SendTimerLog(statMetric, st.remoteservice, key, duration, st.code, tagOptional)

}

func EndStat(st *StatEntry, category string, code int) {

	st.code = code
	st.end = time.Now()
	st.category = category
	ch <- st
	// sendTimeLog(st)
}

// NewStatEntry ....
func NewStatEntry(event string) *StatEntry {
	return &StatEntry{
		event:      event,
		start:      time.Now(),
		splitpoint: ".",
	}
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
