package modules

import (
	"github.com/Sab94/ipfs-monitor/block"
	"github.com/Sab94/ipfs-monitor/client"
	"github.com/Sab94/ipfs-monitor/config"
	bitswapstat "github.com/Sab94/ipfs-monitor/modules/bitswap-stat"
	bootstraplist "github.com/Sab94/ipfs-monitor/modules/bootstrap-list"
	bwstat "github.com/Sab94/ipfs-monitor/modules/bw-stat"
	peerinfo "github.com/Sab94/ipfs-monitor/modules/peer-info"
	repostat "github.com/Sab94/ipfs-monitor/modules/repo-stat"
	swarmpeers "github.com/Sab94/ipfs-monitor/modules/swarm-peers"
	"github.com/rivo/tview"
)

// BootstrapAllModules loads all the modules
func BootstrapAllModules(cfg *config.Config, hc *client.HttpClient,
	app *tview.Application) []block.Block {
	blocks := []block.Block{}

	bootstrap := bootstraplist.NewWidget(cfg, hc, app)
	if bootstrap != nil {
		blocks = append(blocks, bootstrap)
	}

	swarmPeers := swarmpeers.NewWidget(cfg, hc, app)
	if swarmPeers != nil {
		blocks = append(blocks, swarmPeers)
	}

	repoStat := repostat.NewWidget(cfg, hc, app)
	if repoStat != nil {
		blocks = append(blocks, repoStat)
	}

	peerInfo := peerinfo.NewWidget(cfg, hc, app)
	if peerInfo != nil {
		blocks = append(blocks, peerInfo)
	}

	bitswapStat := bitswapstat.NewWidget(cfg, hc, app)
	if bitswapStat != nil {
		blocks = append(blocks, bitswapStat)
	}

	bwStat := bwstat.NewWidget(cfg, hc, app)
	if bwStat != nil {
		blocks = append(blocks, bwStat)
	}

	return blocks
}
