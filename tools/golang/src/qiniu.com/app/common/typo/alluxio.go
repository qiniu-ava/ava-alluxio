package typo

type tierInfo struct {
	TierName string `json:"tierName"`
	Value    string `json:"value"`
}

type workerAddress struct {
	DataPort         int32  `json:"dataPort"`
	WebPort          int32  `json:"webPort"`
	RPCPort          int32  `json:"rpcPort"`
	DomainSocketPath string `json:"domainSocketPath"`
	TieredIdentity   struct {
		Tiers []tierInfo `json:"tiers"`
	} `json:"tieredIdentity"`
	Host string `json:"host"`
}

type blockLocaltionInfo struct {
	WorkerID      int64         `json:"workerId"`
	WorkerAddress workerAddress `json:"workerAddress"`
	TierAlias     string        `json:"tierAlias"`
}

type blockInfo struct {
	BlockInfo struct {
		Locations []blockLocaltionInfo `json:"localtions"`
		BlockID   int64                `json:"blockId"`
		Length    int64                `json:"length"`
	} `json:"blockInfo"`
	UfsLocations []string `json:"ufsLocations"`
	Offset       int64    `json:"offset"`
}

type AlluxioPath struct {
	Group                  string      `json:"group"`
	Folder                 bool        `json:"folder"`
	Cacheable              bool        `json:"cacheable"`
	Mode                   int32       `json:"mode"`
	FileBlockInfos         []blockInfo `json:"fileBlockInfos"`
	FileID                 int64       `json:"fileId"`
	PersistenceState       string      `json:"persistenceState"`
	BlockSizeBytes         int64       `json:"blockSizeBytes"`
	TTL                    int32       `json:"ttl"`
	TTLAction              string      `json:"ttlAction"`
	Persisted              bool        `json:"persisted"`
	MountPoint             bool        `json:"mountPoint"`
	UfsFingerprint         string      `json:"ufsFingerprint"`
	Pinned                 bool        `json:"pinned"`
	CreationTimeMs         int64       `json:"creationTimeMs"`
	LastModificationTimeMs int64       `json:"lastModificationTimeMs"`
	Completed              bool        `json:"completed"`
	MountID                int64       `json:"mountId"`
	BlockIds               []int64     `json:"blockIds"`
	InMemoryPercentage     int32       `json:"inMemoryPercentage"`
	InAlluxioPercentage    int32       `json:"inAlluxioPercentage"`
	UfsPath                string      `json:"ufsPath"`
	Length                 int64       `json:"length"`
	Name                   string      `json:"name"`
	Path                   string      `json:"path"`
	Owner                  string      `json:"owner"`
}

type AlluxioListResult []AlluxioPath
