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
