package service

import (
	"goshop/configs"
	"strconv"
	"time"

	"github.com/golang/glog"
	"github.com/hashicorp/consul/api"
	"gopkg.in/yaml.v3"
)

type ConsulClient struct {
	Conn *api.Client
}

func NewConsulClient(cfg *configs.ConsulConfig) (consulClient *ConsulClient, err error) {
	consulCfg := api.DefaultConfig()
	consulCfg.Address = cfg.Host + ":" + cfg.Port
	consulCfg.Scheme = cfg.Scheme
	conn, err := api.NewClient(consulCfg)
	if err != nil {
		glog.Errorln("[Consul] new consul client error: ", err.Error())
		return nil, err
	}
	return &ConsulClient{
		Conn: conn,
	}, nil
}

func (c *ConsulClient) ServiceRegister(id, name, address, port, checkTimeout, checkInterval string) bool {
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
	err := c.Conn.Agent().ServiceRegister(&reg)
	if err != nil {
		glog.Errorln("[Consul] Service register error: ", err.Error())
		return false
	}
	glog.Infof("[Consul] Service [%s] register [%s, %s] in consul center %s\n",
		reg.Name, reg.ID, reg.Address+":"+port, configs.GetConf().GetConsulAddr())
	return true
}

func (c *ConsulClient) GetAddrByServiceName(serviceName string) (string, error) {
	services, _, err := c.Conn.Health().Service(serviceName, "grpc", true, nil)
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

func (c *ConsulClient) ServiceDeregister(serviceId string) error {
	err := c.Conn.Agent().ServiceDeregister(serviceId)
	if err != nil {
		glog.Errorln("[Consul] Service deregister failed: ", err.Error())
		return err
	}
	return nil
}

func (c *ConsulClient) ServiceRecover(serviceName string) (addr string, err error) {
	var (
		maxRetry      = 5
		retryInterval = 1 * time.Second
	)
	for range maxRetry {
		addr, err = c.GetAddrByServiceName(serviceName)
		if err == nil && addr != "" {
			break
		}
		time.Sleep(retryInterval)
		retryInterval *= 2
	}
	return
}

// Consul configure management.
func (c *ConsulClient) ConfigRecover(configPath string) (*configs.Config, error) {
	var (
		err       error
		cfg       configs.Config
		lastIndex uint64
	)
	kv := c.Conn.KV()
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

func (c *ConsulClient) ConfigQuery(configPath string) (*configs.Config, error) {
	kv := c.Conn.KV()
	data, _, err := kv.Get(configPath, nil)
	var cfg configs.Config
	if err = yaml.Unmarshal(data.Value, &cfg); err != nil {
		glog.Errorln("[Consul] Failed to parse config file: ", err.Error())
		return nil, err
	}
	return &cfg, nil
}
