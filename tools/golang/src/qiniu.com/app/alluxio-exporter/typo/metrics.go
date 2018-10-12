package typo

const (
	AlluxioNamespace = "alluxio"
)

type WorkerStat struct {
	Capacity    capacity               `json:"capacity"`
	StartTimeMs float64                `json:"startTimeMs"`
	CapacityAll workerTierCapacity     `json:"tierCapacity"`
	Metric      map[string]interface{} `json:"metrics"`
}

type capacity struct {
	Total float64 `json:"total"`
	Used  float64 `json:"used"`
}

type workerTierCapacity struct {
	MEM capacity `json:"MEM"`
	SSD capacity `json:"SSD"`
}

type MasterStat struct {
	Capacity    capacity               `json:"capacity"`
	LostWorkers []Worker               `json:"LostWorkers"`
	Workers     []Worker               `json:"workers"`
	CapacityAll masterTierCapacity     `json:"tierCapacity"`
	Metric      map[string]interface{} `json:"metrics"`
	UFSCapacity capacity               `json:"ufsCapacity"`
	StartTimeMs float64                `json:"startTimeMs"`
}

type Worker struct {
	Address        workerAddress `json:"address"`
	Capacity       float64       `json:"capacityBytes"`
	ID             float64       `json:"id"`
	LastContactSec float64       `json:"lastContactSec"`
	StartTimeMs    float64       `json:"startTimeMs"`
	State          string        `json:"state"`
	UsedBytes      float64       `json:"usedBytes"`
}

type workerAddress struct {
	Host     string `json:"host"`
	DataPort int    `json:"dataPort"`
	RPCPort  int    `json:"rpcPort"`
}

type masterTierCapacity struct {
	HDD capacity `json:"HDD"`
	MEM capacity `json:"MEM"`
	SSD capacity `json:"SSD"`
}
