package service

import (
	"goshop/configs"
	"strconv"
	"time"

	"github.com/golang/glog"
	"github.com/hashicorp/consul/api"
	"gopkg.in/yaml.v3"
)

// TODO
// type Consul struct {
// 	client *api.Config
// }

func consulConf() *api.Config {
	cfg := &configs.GetConf().ConsulCfg
	consulConfig := api.DefaultConfig()
	consulConfig.Address = cfg.Host + ":" + cfg.Port
	consulConfig.Scheme = cfg.Scheme // consul api protocol
	return consulConfig
}

func newConsulClient() (consulClient *api.Client, err error) {
	consulClient, err = api.NewClient(consulConf())
	if err != nil {
		glog.Errorln("[Consul] new consul client error: ", err.Error())
	}
	return
}

func ServiceRegister(id, name, address, port, checkTimeout, checkInterval string) bool {
	consulClient, err := newConsulClient()
	if err != nil {
		return false
	}
	pt, _ := strconv.Atoi(port)
	reg := api.AgentServiceRegistration{
		ID:      name + "_" + id,
		Name:    name,
		Tags:    []string{"grpc"},
		Address: address,
		Port:    pt,
		Check: &api.AgentServiceCheck{
			CheckID:  "consul-check " + name + "_" + id,
			TCP:      address + ":" + port,
			Timeout:  checkTimeout,
			Interval: checkInterval,
		},
	}
	err = consulClient.Agent().ServiceRegister(&reg)
	if err != nil {
		glog.Errorln("[Consul] Service register error: ", err.Error())
		return false
	}
	glog.Infof("[Consul] Service [%s] register [%s, %s] in consul center %s\n",
		reg.Name, reg.ID, reg.Address+":"+port, consulConf().Address)
	return true
}

func getAddrByServiceName(serviceName string) (string, error) {
	consulClient, err := newConsulClient()
	if err != nil {
		return "", err
	}
	services, _, err := consulClient.Health().Service(serviceName, "grpc", true, nil)
	if err != nil {
		glog.Errorln("[Consul] Service recover error: ", err.Error())
		return "", err
	}
	if len(services) == 0 {
		glog.Warningf("[Consul] No any available service on [%s].", serviceName)
		return "", nil
	}
	addr := services[0].Service.Address + ":" + strconv.Itoa(services[0].Service.Port)
	return addr, nil
}

func ServiceDeregister(serviceId string) {
	consulClient, err := newConsulClient()
	if err != nil {
		return
	}
	consulClient.Agent().ServiceDeregister(serviceId)
}

func ServiceRecover(serviceName string) (addr string, err error) {
	var (
		maxRetry      = 5
		retryInterval = 1 * time.Second
	)
	for range maxRetry {
		addr, err = getAddrByServiceName(serviceName)
		if err == nil && addr != "" {
			break
		}
		time.Sleep(retryInterval)
		retryInterval *= 2
	}
	return
}

// Consul configure management.
func ConfigRecover(configPath string) (*configs.Config, error) {
	var (
		err       error
		cfg       configs.Config
		lastIndex uint64
	)
	consulClient, err := newConsulClient()
	if err != nil {
		return nil, err
	}
	kv := consulClient.KV()
	data, meta, err := kv.Get(configPath, &api.QueryOptions{
		WaitIndex: lastIndex,
		WaitTime:  600 * time.Second,
	})
	if err != nil {
		glog.Errorf("[Consul] Failed to recover [%s] config file: %s\n", configPath, err.Error())
		return nil, err
	}
	if meta.LastIndex != lastIndex {
		lastIndex = meta.LastIndex
		if data != nil {
			if err = yaml.Unmarshal(data.Value, &cfg); err != nil {
				glog.Errorf("[Consul] Failed to parse [%s] config file: %s\n", configPath, err.Error())
				return nil, err
			}
		}
	}
	return &cfg, nil
}

func ConfigQuery(configPath string) (*configs.Config, error) {
	consulClient, err := newConsulClient()
	if err != nil {
		return nil, err
	}
	kv := consulClient.KV()
	data, _, err := kv.Get(configPath, nil)
	var cfg configs.Config
	if err = yaml.Unmarshal(data.Value, &cfg); err != nil {
		glog.Errorln("[Consul] Failed to parse config file: ", err.Error())
		return nil, err
	}
	return &cfg, nil
}
