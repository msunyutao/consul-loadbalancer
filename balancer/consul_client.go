package balancer

import (
	"errors"
	"strconv"

	"github.com/hashicorp/consul/api"
)

type ConsulClient struct {
	client *api.Client
}

func NewConsulClient(address string) (*ConsulClient, error) {
	config := api.DefaultConfig()
	config.Address = address
	client, err := api.NewClient(config)
	if err != nil {
		return nil, err
	}
	return &ConsulClient{client: client}, nil
}

func (c *ConsulClient) Get(key string) ([]byte, error) {
	res, _, err := c.client.KV().Get(key, nil)
	if err != nil {
		return nil, err
	}
	return res.Value, nil
}

func (c *ConsulClient) Put(key string, val []byte) error {
	_, err := c.client.KV().Put(&api.KVPair{Key: key, Value: val}, nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *ConsulClient) GetServiceNodes(service, tag string, passingOnly bool) ([]ServiceNode, error) {
	entries, _, err := c.client.Health().Service(service, tag, passingOnly, nil)
	if err != nil {
		return nil, err
	}
	if len(entries) == 0 {
		return nil, errors.New("consul has not that service or tag")
	}
	serviceNodes := make([]ServiceNode, len(entries))
	for i, entry := range entries {
		serviceNode := ServiceNode{}
		serviceNode.Zone = entry.Service.Meta["zone"]
		balanceFactor, _ := strconv.ParseFloat(entry.Service.Meta["balanceFactor"], 64)
		serviceNode.BalanceFactor = balanceFactor
		serviceNode.InstanceID = entry.Service.Meta["instanceID"]
		serviceNode.PublicIP = entry.Service.Meta["publicIP"]
		serviceNode.Host = entry.Service.Address
		serviceNode.Port = entry.Service.Port
		serviceNodes[i] = serviceNode
	}
	return serviceNodes, nil
}
