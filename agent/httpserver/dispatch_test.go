package httpserver

import (
	"encoding/json"
	"github.com/chaosblade-io/chaosblade/agent/config"
	"testing"
	"time"
)

func TestWorkflow_Parse(t *testing.T) {
	config.ConfInit = func() {}

	//	  "sys_time": 1566040202000,
	jsonByte := []byte(`{
 "job_id": "J1",
 "portal_time": 1567825776000,
 "flow_starttime":1567825786000,
 "events": [
   {
     "event_id": "E1",
     "command": "cpu",
     "sub_command": "fullload",
     "flag": "cpu-count=2",
     "start_offset": 30000,
     "end_offset": 60000
   },{
     "event_id": "E2",
     "command": "cpu",
     "sub_command": "fullload",
     "flag": "cpu-count=2",
     "start_offset": 600000,
     "end_offset":900000
   },{
     "event_id": "E3",
     "command": "jvm",
     "sub_command": "throwCustomException",
     "flag": "classname=com.icbc.xxx exception=java.lang.Exception",
     "start_offset": 900000,
     "end_offset":1200000
   }
 ],
"prepare":[
{
     "prepare_id": "E3",
     "command": "prepare",
     "sub_command": "jvm",
     "flag": "process=TestProvider"
}
]
}`)
	wf := &workflow{}
	err := json.Unmarshal(jsonByte, wf)
	if err != nil {
		t.Error(err)
	}
	nowTime := time.Now().Truncate(time.Millisecond)
	t.Log(nowTime)
	wf.Parse(nowTime)
	result, err := json.Marshal(wf)
	if err != nil {
		t.Error(err)
	}
	nodetimeGap := nowTime.UnixNano()/int64(time.Millisecond) - wf.PortalTime

	startTime := nowTime.Add(40000 * time.Millisecond)
	endTime := nowTime.Add(70000 * time.Millisecond)

	flag := false
	for i, _ := range wf.Events {
		if wf.Events[i].ExpId == "E1" {
			flag = true
			if wf.Events[i].JobId != "J1" {
				t.Error("not J1")
			}
			if wf.Events[i].NodetimeGap != nodetimeGap {
				t.Error("NodetimeGap incorrect.")
			}
			if wf.Events[i].Status != "0" {
				t.Error("J1.E3.status!=0")
			}
			if wf.Events[i].StartTime.UnixNano() != startTime.UnixNano() {
				t.Error("J1.E3.StartTime not right :", wf.Events[i].StartTime, startTime)
			}
			if wf.Events[i].EndTime.UnixNano() != endTime.UnixNano() {
				t.Error("J1.E3.EndTime not right :", wf.Events[i].EndTime, endTime)
			}
		}
	}
	if !flag {
		t.Error("there is no E1")
	}
	t.Log(string(result))
}
