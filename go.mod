module github.com/chaosblade-io/chaosblade

go 1.13

require (
	github.com/Nvveen/Gotty v0.0.0-20120604004816-cd527374f1e5 // indirect
	github.com/chaosblade-io/chaosblade-exec-docker v0.5.0
	github.com/chaosblade-io/chaosblade-exec-os v0.5.0
	github.com/chaosblade-io/chaosblade-operator v0.5.0
	github.com/chaosblade-io/chaosblade-spec-go v0.5.0
	github.com/coreos/etcd v3.3.22+incompatible // indirect
	github.com/dustin/go-humanize v1.0.0 // indirect
	github.com/etcd-io/etcd v3.3.22+incompatible
	github.com/fsnotify/fsnotify v1.4.7
	github.com/mattn/go-sqlite3 v1.10.1-0.20190217174029-ad30583d8387
	github.com/prometheus/client_golang v1.6.0
	github.com/prometheus/common v0.10.0
	github.com/prometheus/node_exporter v1.0.0
	github.com/shirou/gopsutil v2.19.9+incompatible
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/cobra v0.0.4-0.20190109003409-7547e83b2d85
	github.com/spf13/pflag v1.0.4-0.20181223182923-24fa6976df40
	github.com/spf13/viper v1.3.2
	golang.org/x/crypto v0.0.0-20200510223506-06a226fb4e37
	gopkg.in/alecthomas/kingpin.v2 v2.2.6
	k8s.io/apimachinery v0.17.0
	k8s.io/client-go v11.0.0+incompatible
)

// Pinned to kubernetes-1.13.11
replace (
	bitbucket.org/ww/goautoneg => github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822
	github.com/chaosblade-io/chaosblade-exec-os => github.com/yixy/chaosblade-exec-os v0.5.1
	github.com/prometheus/client_golang => github.com/prometheus/client_golang v0.9.3-0.20190127221311-3c4408c8b829
	github.com/prometheus/common => github.com/prometheus/common v0.3.0
	github.com/prometheus/node_exporter => github.com/prometheus/node_exporter v0.18.1
	k8s.io/api => k8s.io/api v0.0.0-20190817221950-ebce17126a01
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.0.0-20190919022157-e8460a76b3ad
	k8s.io/apimachinery => k8s.io/apimachinery v0.0.0-20190817221809-bf4de9df677c
	k8s.io/client-go => k8s.io/client-go v0.0.0-20190817222206-ee6c071a42cf
)
