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

func BootstrapAllModules(cfg *config.Config, hc *client.HttpClient,
	app *tview.Application) []block.Block {
	blocks := []block.Block{}

	bootstrap := bootstrap_list.NewWidget(cfg, hc, app)
	blocks = append(blocks, bootstrap)

	swarmPeers := swarm_peers.NewWidget(cfg, hc, app)
	blocks = append(blocks, swarmPeers)

	repoStat := repo_stat.NewWidget(cfg, hc, app)
	blocks = append(blocks, repoStat)

	peerInfo := peer_info.NewWidget(cfg, hc, app)
	blocks = append(blocks, peerInfo)

	bitswapStat := bitswap_stat.NewWidget(cfg, hc, app)
	blocks = append(blocks, bitswapStat)

	bwStat := bw_stat.NewWidget(cfg, hc, app)
	blocks = append(blocks, bwStat)

	return blocks
}
