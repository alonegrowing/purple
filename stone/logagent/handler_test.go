package logagent

import (
	"testing"
)

var agentclient *LogAgent

func init() {
	agentclient = InitAgent()
}

func TestTimerLogMsg(t *testing.T) {

	serviceName := "longlink.push"   //本实例的服务名称
	remoteService := "live.liveinfo" // 调用远端服务
	metric := "getLiveInfo"          // 监控的信息,
	duration := 1000                 // 耗时
	code = 0                         // 错误码
	agentclient.SendTimerLogMsg(statMetric, st.remoteservice, metric, duration, code)
}

func TestTimerLog(t *testing.T) {
	serviceName := "longlink.push"   //本实例的服务名称
	remoteService := "live.liveinfo" // 调用远端服务
	metric := "getLiveInfo"          // 监控的信息,
	duration := 1000                 // 耗时
	code = 0                         // 错误码
	lat := NewLogAgentTagOptions()
	lat.SetTag("key", "value")
	agentclient.SendTimerLog(statMetric, st.remoteservice, metric, duration, code, lat)
}

func TestSnapshot(t *testing.T) {

	serviceName := "longlink.push" // 本机实例的服务名称
	snapshotValue := 100           //  快照数量
	lat := NewLogAgentTagOptions()
	lat.SetMemCache("liveidcount") // 设置当前是统计内存数量的tag 模式

	// lat.SetTag("key", "value") //自定义tag类型
	agentclient.SendSnapshotLog(service, snapshotValue, lat)
}
