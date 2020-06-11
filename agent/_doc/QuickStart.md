# quick start #

## chaos-agent ##

chaos-agent启动如下：(系统管理员权限，如root)，详细命令参数请参考`./blade httpagent -h`

```
#for test mode, no etcd regist
./blade httpagent
#for virtual (虚拟机）,with etcd regist：/chaos/virtual/F-APIP/hostIp/hostPort
./blade httpagent --environment=1 --host-ip=122.66.29.88 --host-port=8081 --app-shortname=F-APIP
#for container（容器）,with etcd regist：/chaos/container/F-APIP/hostIp/hostPort/containerId/containerPort
./blade httpagent --environment=2 --host-ip=122.66.29.88 --host-port=8081 --container-id=1234234 --container-port=36661 --app-shortname=F-APIP
```

chaos实验支持场景如下：

```
场景变更历史：

[chaos-agent-0.1.1/chaosblade-0.3.0]
* 磁盘占用冲高：--mount-point修改为--path
* 磁盘IO冲高：--mount-point修改为--path
* 磁盘IO冲高：去除--count (count原本用于指定IO模拟操作每次读写的块数目，0.3.0版本中chaosblade将count写死为100常量，不用外部传入了)
* 网络延迟：增加--destination-ip
* 网络丢包：增加--destination-ip

[chaos-agent-0.1.3/chaosblade-0.4.0]
* 磁盘占用冲高： 增加--percent指定百分比
```

| scene        | command | subcommand | flag(separate by blank)                                      | flag comment                                                 | cli example                                                  | chaos binary file(演练可执行程序） |
| ------------ | ------- | ---------- | ------------------------------------------------------------ | ------------------------------------------------------------ | ------------------------------------------------------------ | ---------------------------------- |
| CPU冲高      | cpu     | fullload   | --cpu-percent=20 --cpu-count=1                               | cpu-percent指定CPU使用率，1-100的整数，必输; cpu-count指定CPU核数，正整数（大于实际核数会报错），非必输，注意 | blade create cpu fullload --cpu-percent=20； blade create cpu fullload --cpu-percent=20 cpu-count=1 | chaos_burncpu                      |
| 内存冲高     | mem     | load       | --mem-percent                                                | mem-percent指定内存使用率，1-100的整数，必输                 | blade create mem load --mem-percent=10                       | chaos_burnmem                      |
| 磁盘占用冲高 | disk    | fill       | --size=1024 --path=/ --percent=50      | percent指定磁盘文件系统填充百分比，若percent和size同时设置则percent优先生效，正整数；size指定磁盘文件系统填充大小，默认MB，整数； path指定需要填充的文件系统挂载点，目录不存在或不是磁盘挂载点会报错，必输 | blade create   --size=1024 --path=/                          | chaos_filldisk                     |
| 磁盘IO冲高   | disk    | burn       | --write=true --read=true --path=/ --size=1                | write指定写IO；read指定读IO；（write和read必须至少有一个是true）path指定需要模拟IO冲高的文件系统的挂载点，目录不存在或不是磁盘挂载点会报错，必输；size指定IO模拟操作每次读写的块的大小，必输，单位MB（注意，read场景会创建一个600MB的文件用于作为read的源，并在实验完成后销毁这个文件）。 | blade create disk burn  --write=true --read=false --path=/  --size=1 | chaos_burnio                       |
| 网络延迟     | network | delay      | --interface=eth0 --exclude-port=22 --local-port=80 --remote-port=10099 --time=100 --offset=20 --destination-ip | local-port指定针对访问哪个本地端口进行延迟注入；remote-port指定针对访问远端哪个端口进行延迟注入；exclude-port指定全场景网络延迟注入场景下的例外端口。【local-port、remote-port、exclude-port不能全为空；并且当local-port或remote-port被指定时，exclude-port是非法的；】time指定延时时间，单位ms，必输；offset指定延时抖动，增加一个随机时间长度，让延迟时间出现在某个范围，如果设置为 20ms，那么报文延迟的时间在 time ± 20ms 之间（以time=100ms，offset=20ms为例，实际延时范围为90ms - 110ms），单位ms，必输；destination-ip用于指定远端目标IP，可以使用192.168.1.1或者192.168.1.1/32两种格式来指定，非必输 | blade create network delay --interface=eth0 --local-port=8080 --time=5000 --offset=1000; blade create network delay --interface=eth0 --local-port=8080 --remote-port=80 --time=5000 --offset=1000; blade create network delay --interface=eth0 --exclude-port=22  --time=5000 --offset=1000 | chaos_dlnetwork                    |
| 网络丢包     | network | loss       | --interface=eth0 --exclude-port=22 --local-port=80 --remote-port=10099 --percent=50 --destination-ip | percent指定丢包率，为0-100的整数，必输。其余flag与网络延迟场景含义相似 | blade create network loss --interface=eth0 --exclude-port=22  --percent=50 | chaos_dlnetwork                    |
| 端口封禁     | network | drop       | --local-port=80 --remote-port=10099                          | local-port指定针对访问哪个本地端口进行封禁；remote-port指定针对访问远端哪个端口进行封禁；与网络延迟场景含义相似 | blade create network drop --interface=eth0 --local-port=80   | chaos_dropnetwork                  |
| DNS解析      | network | dns        | --ip=122.20.25.1  --domain=www.github.com                    | 将domain指定域名解析至ip，必输。                             | blade create network dns --domain=www.baidu.com --ip=172.20.10.3 | chaos_changedns                    |
| 待支持       | jvm     |            |                                                              |                                                              |                                                              |                                    |
| 待支持       | dubbo   |            |                                                              |                                                              |                                                              |                                    |
| 待支持       | process |            |                                                              |                                                              |                                                              |                                    |
| 待支持       | http    |            |                                                              |                                                              |                                                              |                                    |
| 待支持       | docker  |            |                                                              |                                                              |                                                              |                                    |
| 待支持       | k8s     |            |                                                              |                                                              |                                                              |                                    |
| 待支持       | mysql   |            |                                                              |                                                              |                                                              |                                    |

## api/exp/dispatch ##

任务下发接口，用于发布实验注入计划。

> 关于接口的说明：编排演练计划时，需要保证同一时刻，不能有多个同类型的演练同时执行。这里同类型的演练，按照chaos binary file(演练可执行程序）来区分。比如，网络延迟和网络丢包两个场景的演练，都对应chaos_dlnetwork程序，是同一个类型的演练，不能同时执行。又如，CPU冲高和内存冲高，分别对应chaos_burncpu和chaos_burnmem两个程序，是不同类型的演练，可以同时执行。

** 请求参数 **

通讯区字段定义：

| 字段名称       | 字段类型              | 字段长度 | 是否必输 | 说明                                                         |
| -------------- | --------------------- | -------- | -------- | ------------------------------------------------------------ |
| job_id         | string                | 20       | 是       | chaos-job定义的job_id，格式为flowid\|taskid\|ip\|port        |
| portal_time    | int64                 |          | 是       | 混沌工程管理平台的系统时间，传值按Unix-time格式换算成毫秒。Unix-time格式指seconds since 1970-01-01 00:00:00 UTC，可以使用`date +s`命令查看。例如，Tue Sep 17 08:53:30 CST 2019 对应的Unix-time为1568681610，转换为portal_time则是1568681610*1000=1568681610000 |
| flow_starttime | int64                 |          | 是       | 演练流程开始时间，传值按Unix-time格式换算成毫秒。            |
| duration | int64                 |          | 是       | 演练流程持续时间，毫秒。            |
| events         | []ExperimentPlanModel |          | 是       | envents是由event元素组成的数组。每个event对应一个ExperimentPlanModel结构体 |

event元素（ExperimentPlanModel结构体）定义：

| 字段名称     | 字段类型 | 字段长度 | 是否必输 | 说明                                                         |
| ------------ | -------- | -------- | -------- | ------------------------------------------------------------ |
| event_id     | string   | 32       | 是       | 用于唯一标识一个下发的实验，和实验真实执行时生成的uid一一对应。 |
| command      | string   | -        | 是       | chaos实验场景中的command                                     |
| sub_command  | string   | -        | 是       | chaos实验场景中的sub_command                                 |
| flag         | string   | -        | 是       | chaos实验场景中的flag                                        |
| start_offset | int64    |          | 是       | 实验开始时间偏移量（实验计划开始时间-flow_starttime的差值），传值按Unix-time格式换算成毫秒。 |
| end_offset   | int64    |          | 是       | 实验结束时间偏移量（实验计划结束时间-flow_starttime的差值），传值按Unix-time格式换算成毫秒。 |

** 响应参数 **

| 字段名称 | 字段类型 |
| -------- | -------- |
| code     | int64    |
| success  | bool     |
| result   |          |

** 报文示例 **

该接口是一个HTTP接口，接口请求示例报文如下：

```
POST http://122.66.29.88:36661/chaosagent/api/exp/dispatch
Content-Type: application/json

{
  "job_id": "flowid|taskid|ip|port",
  "portal_time": 1567825776000,
  "flow_starttime":1567825776000,  
  "duration":3600000,
  "events": [
    {
      "event_id": "E1",
      "command": "cpu",
      "sub_command": "fullload",
      "flag": " cpu-count=1 cpu-percent=20",
      "start_offset": 10000,
      "end_offset":20000
    },{
      "event_id": "E2",
      "command": "cpu",
      "sub_command": "fullload",
      "flag": "cpu-count=2  cpu-percent=20 ",
      "start_offset": 20000,
      "end_offset":30000
    },{
      "event_id": "E3",
      "command": "cpu",
      "sub_command": "fullload",
      "flag": " cpu-count=3 cpu-percent=20 ",
      "start_offset": 30000,
      "end_offset":40000
    }
  ]
}
```

响应

```
{
	"code": 200,
	"success": true,
	"result": "ok"
}
```


## api/exp/terminate ##

终止演练流程

** 请求参数 **

无

** 响应参数 **

| 字段名称 | 字段类型 |
| -------- | -------- |
| code     | int64    |
| success  | bool     |
| result   |          |

** 报文示例 **

请求

```
POST http://122.66.29.88:36661/chaosagent/api/exp/terminate
```

响应

```
{
	"code": 200,
	"success": true,
	"result": "Terminated exp(event_id|uid) list is [E333| E222|67cbeee3abc61ffa] , UnTerminated exp(event_id|uid) list is []"
}
```

## api/metrics ##

** 请求参数 **

| 字段名称 | 字段类型 | 字段长度 | 是否必输 | 说明                                                         |
| -------- | -------- | -------- | -------- | ------------------------------------------------------------ |
| job_id   | string   | 20       | 是       | 注意：如果不输入或输入job_id不存在，则仅查询节点metrics信息，不查询具体演练events信息 |

** 响应参数 **

| 字段名称     | 字段类型 | 字段长度 | 是否必输 | 说明                                                         |
| ------------ | -------- | -------- | -------- | ------------------------------------------------------------ |
| job_id       | string   | 20       |          |                                        |
| nodetime_gap | int64    |          | 是       | 混沌工程平台与演练服务器的实际差值。node_sys_time-portal_time，其中portal_time为任务下发接口传值，node_sys_time为演练目标服务器系统时间。格式均为Unix-time转换为毫秒。 |
| finished       |    string     |    1    |    是   | 演练流程状态信息，0-流程待执行（当前时间小于流程开始时间） 1-流程执行中（当前时间在流程开始和结束时间之间） 2-流程执行完（当前时间在流程结束时间之后）       |
| events       |          |          |          | 演练event状态信息，若job_id不存在则返回空                    |
| metrics      |          |          | 是       | 节点系统资源指标                                             |

events

| 字段名称     | 字段类型 | 字段长度 | 是否必输 | 说明                                                       |
| ------------ | -------- | -------- | -------- | ---------------------------------------------------------- |
| event_id     | string   | 32       | 是       |                                                            |
| uid          | string   | 32       | 是       | chaosblade的uid                                            |
| event_status | string   | 1        | 是       | 0-待执行，1-演练执行中，2-执行成功，3-执行失败，4-演练终止 |
| err_msg      | string   |          | 否       |                                                            |

metrics

| 字段名称                        | 字段类型 | 说明                                                  |
| ------------------------------- | -------- | ----------------------------------------------------- |
| sample_time                     | string   | 采样时间，取节点采样时的系统时间，Unix-time转换为毫秒 |
| sample_interval                     | string   | 采样间隔时间周期，默认1s |
| load1                           | string   | 过去1min节点负载情况                                  |
| load5                           | string   | 过去5min节点负载情况                                  |
| load15                          | string   | 过去15min节点负载情况                                 |
| cpu_seconds_total               | string   | 本次开机以来CPU的总耗时，秒数                         |
| cpu_seconds_idle                | string   | 本次开机以来CPU的空闲耗时，秒数                       |
| cpu_seconds_total_ls            | string   | 上次采样的cpu_seconds_total                           |
| cpu_seconds_idle_ls             | string   | 上次采样的cpu_seconds_idle                            |
| cpu_usage                       | string   | 两次采样间隔内CPU使用率                               |
| memory_memtotal_bytes           | string   | 内存总大小，字节                                      |
| memory_mem_buffers_bytes        | string   | 内存buffer大小，字节                                  |
| memory_cached_bytes             | string   | 内存cache大小，字节                                   |
| memory_mem_free_bytes           | string   | 内存free大小，字节                                    |
| memory_usage                    | string   | 内存使用率                                            |
| filesystem_size_bytes           | string   | 文件系统大小，字节                                    |
| filesystem_free_bytes           | string   | 文件系统空闲大小，字节                                |
| filesystem_usage                | string   | 文件系统使用率                                        |
| filefd_allocated                | string   | 已分配的打开文件描述符数（文件句柄）                  |
| filefd_maximun                  | string   | 系统的打开文件描述符数上限                            |
| sockstat_tcp_alloc              | string   | tcp连接数                                             |
| sockstat_tcp_tw                 | string   | 处于time_wait状态的tcp连接数                          |
| network_receive_bytes_total     | string   | 本次开机以来接收到的网络接收报文大小，字节            |
| network_transmit_bytes_total    | string   | 本次开机以来发送到的网络请求报文大小，字节            |
| network_receive_bytes_total_ls  | string   | 上次采样的network_receive_bytes_total                 |
| network_transmit_bytes_total_ls | string   | 上次采样的network_transmit_bytes_total                |
| network_receive_bytes           | string   | 两次采样间隔内的网络接收报文大小，字节                |
| network_transmit_bytes          | string   | 两次采样间隔内的网络请求报文大小，字节                |
| disk_read_bytes_total          | string   | 本次开机以来磁盘读取字节数             |
| disk_write_bytes_total          | string   | 本次开机以来磁盘写入字节数                |
| disk_read_bytes_total_ls          | string   | 上次采样的disk_read_bytes_total                |
| disk_write_bytes_total_ls          | string   | 上次采样的disk_write_bytes_total                |
| disk_read_bytes          | string   | 两次采样间隔内的磁盘读取字节数                |
| disk_write_bytes          | string   | 两次采样间隔内的磁盘写入字节数               |

** 报文示例 **

请求

```
GET http://122.66.29.88:36661/chaosagent/api/metrics?job_id=flowid|taskid|ip|port
```

响应

```
{
	"job_id": "flowid|taskid|ip|port",
	"nodetime_gap": 890440937,
    "finished":false,
	"events": [
		{
			"event_id": "E1",
			"event_status": "2",
			"err_msg": ""
		},
		{
			"event_id": "E2",
			"event_status": "1",
			"err_msg": ""
		},
		{
			"event_id": "E3",
			"event_status": "0",
			"err_msg": ""
		}
	],
	"metrics": {
		"sample_time": 1570501539362,
		"sample_interval": 1,
		"load1": "1.49",
		"load5": "0.89",
		"load15": "0.35",
		"cpu_seconds_total": "23116008.610000",
		"cpu_seconds_idle": "22988368.590000",
		"cpu_seconds_total_ls": "23116004.560000",
		"cpu_seconds_idle_ls": "22988364.590000",
		"cpu_usage": "0.012346",
		"memory_memtotal_bytes": "7978545152",
		"memory_mem_buffers_bytes": "379346944",
		"memory_cached_bytes": "3749363712",
		"memory_mem_free_bytes": "2545491968",
		"memory_usage": "0.163481",
		"filesystem_size_bytes": "133885157376",
		"filesystem_free_bytes": "120961809408",
		"filesystem_usage": "0.096526",
		"filefd_allocated": "1472",
		"filefd_maximun": "6815744",
		"sockstat_tcp_alloc": "20",
		"sockstat_tcp_tw": "0",
		"network_receive_bytes_total": "25994872951",
		"network_transmit_bytes_total": "21454807242",
		"network_receive_bytes_total_ls": "25994872951",
		"network_transmit_bytes_total_ls": "21454807242",
		"network_receive_bytes": "0",
		"network_transmit_bytes": "0",
		"disk_read_bytes_total": "108016273408",
		"disk_write_bytes_total": "307076105216",
		"disk_read_bytes_total_ls": "108016273408",
		"disk_write_bytes_total_ls": "307076105216",
		"disk_read_bytes": "0",
		"disk_write_bytes": "0"
	}
}
```

## api/report ##

演练结果报告查询

** 报文示例 **

请求

```
GET http://122.66.29.88:36661/chaosagent/api/report
```

响应是一个zip包文件流，解压zip包后是演练报告的CSV文件，文件内容如下

```
sample_time,sample_interval,load1,load5,load15,cpu_seconds_total,cpu_seconds_idle,cpu_seconds_total_ls,cpu_seconds_idle_ls,cpu_usage,memory_memtotal_bytes,memory_mem_buffers_bytes,memory_cached_bytes,memory_mem_free_bytes,memory_usage,filesystem_size_bytes,filesystem_free_bytes,filesystem_usage,filefd_allocated,filefd_maximun,sockstat_tcp_alloc,sockstat_tcp_tw,network_receive_bytes_total,network_transmit_bytes_total,network_receive_bytes_total_ls,network_transmit_bytes_total_ls,network_receive_bytes,network_transmit_bytes,disk_read_bytes_total,disk_write_bytes_total,disk_read_bytes_total_ls,disk_write_bytes_total_ls,disk_read_bytes,disk_write_bytes
1570501535327,1,1.62,0.9,0.35,23115992.430000,22988352.600000,,,0.005522,7978545152,379346944,3749130240,2545971200,0.163450,133885157376,120961809408,0.096526,1472,6815744,20,0,25994872808,21454807086,,,25994872808,21454807086,108016273408,307076105216,,,108016273408,307076105216
```