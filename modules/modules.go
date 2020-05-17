package modules

import (
	"github.com/Sab94/ipfs-monitor/block"
	"github.com/Sab94/ipfs-monitor/client"
	"github.com/Sab94/ipfs-monitor/config"
	bitswap_stat "github.com/Sab94/ipfs-monitor/modules/bitswap-stat"
	bootstrap_list "github.com/Sab94/ipfs-monitor/modules/bootstrap-list"
	bw_stat "github.com/Sab94/ipfs-monitor/modules/bw-stat"
	peer_info "github.com/Sab94/ipfs-monitor/modules/peer-info"
	repo_stat "github.com/Sab94/ipfs-monitor/modules/repo-stat"
	swarm_peers "github.com/Sab94/ipfs-monitor/modules/swarm-peers"
	"github.com/rivo/tview"
)

// BootstrapAllModules loads all the modules
func BootstrapAllModules(cfg *config.Config, hc *client.HttpClient,
	app *tview.Application) []block.Block {
	blocks := []block.Block{}

	bootstrap := bootstrap_list.NewWidget(cfg, hc, app)
	if bootstrap != nil {
		blocks = append(blocks, bootstrap)
	}

	swarmPeers := swarm_peers.NewWidget(cfg, hc, app)
	if swarmPeers != nil {
		blocks = append(blocks, swarmPeers)
	}

	repoStat := repo_stat.NewWidget(cfg, hc, app)
	if repoStat != nil {
		blocks = append(blocks, repoStat)
	}

	peerInfo := peer_info.NewWidget(cfg, hc, app)
	if peerInfo != nil {
		blocks = append(blocks, peerInfo)
	}

	bitswapStat := bitswap_stat.NewWidget(cfg, hc, app)
	if bitswapStat != nil {
		blocks = append(blocks, bitswapStat)
	}

	bwStat := bw_stat.NewWidget(cfg, hc, app)
	if bwStat != nil {
		blocks = append(blocks, bwStat)
	}

	return blocks
}
