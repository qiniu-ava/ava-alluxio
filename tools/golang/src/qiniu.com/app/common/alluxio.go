package common

type ZooKeeper struct {
	Enabled  bool     `json:"enabled"`
	Servers  []string `json:"servers"`
	Leader   string   `json:"leader"`
	Election string   `json:"election"`
}

type AlluxioConfig struct {
	MasterHost string    `json:"masterHost,omitempty"`
	ZooKeeper  ZooKeeper `json:"zookeeper,omitempty"`
}
