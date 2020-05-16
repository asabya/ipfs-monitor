package modules

import (
	"github.com/Sab94/ipfs-monitor/block"
	"github.com/Sab94/ipfs-monitor/client"
	"github.com/Sab94/ipfs-monitor/config"
	bootstrap_list "github.com/Sab94/ipfs-monitor/modules/bootstrap-list"
	peer_info "github.com/Sab94/ipfs-monitor/modules/peer-info"
	repo_stat "github.com/Sab94/ipfs-monitor/modules/repo-stat"
)

func BootstrapAllModules(cfg *config.Config, hc *client.HttpClient) []block.Block {
	blocks := []block.Block{}

	bootstrap := bootstrap_list.NewWidget(cfg, hc)
	blocks = append(blocks, bootstrap)

	repoStat := repo_stat.NewWidget(cfg, hc)
	blocks = append(blocks, repoStat)

	peerInfo := peer_info.NewWidget(cfg, hc)
	blocks = append(blocks, peerInfo)

	return blocks
}
