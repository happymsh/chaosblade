package etcd

import (
	"context"
	"github.com/etcd-io/etcd/clientv3"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/chaosblade-io/chaosblade/agent/config"
)

var Client *clientv3.Client
var Lease clientv3.Lease
var LeaseID clientv3.LeaseID
var LeaseRespChan <-chan *clientv3.LeaseKeepAliveResponse

func ClientInit(ctx context.Context) (err error) {
	Client, err = clientv3.New(clientv3.Config{
		Endpoints:         config.EtcdEndpoints,
		DialTimeout:       config.EtcdDialTimeout,
		DialKeepAliveTime: config.EtcdDialKeepAliveTime,
	})
	if err != nil {
		logrus.Error("etcd init error:", err.Error())
		return err
	}

	Lease = clientv3.NewLease(Client)
	leaseResp, err := Lease.Grant(ctx, config.EtcdBeatTime)
	if err != nil {
		logrus.Error("Lease.Grant error:", err.Error())
		return err
	}
	LeaseID = leaseResp.ID
	LeaseRespChan, err = Lease.KeepAlive(ctx, LeaseID)
	if err != nil {
		logrus.Error("Lease.KeepAlive error:", err.Error())
		return err
	}
	return nil
}

func Available(ctx context.Context, key string) (err error) {
	logrus.Info("get oldKey from config")
	kv := clientv3.NewKV(Client)
	oldKey := viper.GetString("etcd-key")
	logrus.Info("oldKey is :", oldKey)
	if oldKey != "" {
		logrus.Info("delete etcd's oldKey.")
		_, err = kv.Delete(ctx, oldKey)
		if err != nil {
			logrus.Error("delete etcd's oldkey failed.", err.Error())
			return err
		}
	}
	logrus.Info("set etcd-key into config:", key)
	viper.Set("etcd-key", key)
	err = viper.WriteConfig()
	if err != nil {
		logrus.Error("write etcd-key into config failed.", err.Error())
		return err
	}
	logrus.Info("put into etcd :", key)
	_, err = kv.Put(ctx, key, "ok", clientv3.WithLease(LeaseID))
	if err != nil {
		logrus.Error("put etcd failed.", err.Error())
		return err
	}

	return err
}
