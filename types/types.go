package types

type Identity struct {
	Addresses       []string `json:"Addresses"`
	AgentVersion    string   `json:"AgentVersion"`
	ID              string   `json:"ID"`
	ProtocolVersion string   `json:"ProtocolVersion"`
	PublicKey       string   `json:"PublicKey"`
}

type RepoStat struct {
	NumObjects uint64 `json:"NumObjects"`
	RepoPath   string `json:"RepoPath"`
	RepoSize   uint64 `json:"RepoSize"`
	StorageMax uint64 `json:"StorageMax"`
	Version    string `json:"Version"`
}

type BootstrapList struct {
	Peers []string `json:"Peers"`
}

type SwarmPeers struct {
	Peers []struct {
		Addr      string `json:"Addr"`
		Direction int    `json:"Direction"`
		Latency   string `json:"Latency"`
		Muxer     string `json:"Muxer"`
		Peer      string `json:"Peer"`
		Streams   []struct {
			Protocol string `json:"Protocol"`
		} `json:"Streams"`
	} `json:"Peers"`
}

type BitswapStat struct {
	BlocksReceived   uint64        `json:"BlocksReceived"`
	BlocksSent       uint64        `json:"BlocksSent"`
	DataReceived     uint64        `json:"DataReceived"`
	DataSent         uint64        `json:"DataSent"`
	DupBlksReceived  uint64        `json:"DupBlksReceived"`
	DupDataReceived  uint64        `json:"DupDataReceived"`
	MessagesReceived uint64        `json:"MessagesReceived"`
	Peers            []interface{} `json:"Peers"`
	ProvideBufLen    int           `json:"ProvideBufLen"`
	Wantlist         []interface{} `json:"Wantlist"`
}

type BWStat struct {
	RateIn   float64 `json:"RateIn"`
	RateOut  float64 `json:"RateOut"`
	TotalIn  uint64   `json:"TotalIn"`
	TotalOut uint64   `json:"TotalOut"`
}
