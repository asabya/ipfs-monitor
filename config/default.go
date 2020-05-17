package config

const defaultConfigFile = `monitor:
  colors:
    border:
      focusable: darkslateblue
      focused: orange
      normal: gray
    background: black
    text: white
  grid:
    columns: [28, 28, 28, 28, 90]
    rows: [4, 10, 10, 10, 4, 90]
    background: black
    border: false
  refreshInterval: 1
  widgets:
    peerinfo:
      enabled: true
      position:
        top: 0
        left: 0
        height: 1
        width: 2
      refreshInterval: 0
      title: "Peer Info"
    bootstraplist:
      enabled: true
      position:
        top: 1
        left: 0
        height: 2
        width: 3
      refreshInterval: 0
      title: "Bootstraps"
    repostat:
      enabled: true
      position:
        top: 3
        left: 0
        height: 1
        width: 2
      refreshInterval: 500
      title: "Repo Stats"
    swarmpeers:
      enabled: true
      position:
        top: 1
        left: 3
        height: 2
        width: 3
      refreshInterval: 10
      title: "Swarm Peers"
    bitswapstat:
      enabled: true
      position:
        top: 0
        left: 4
        height: 1
        width: 2
      refreshInterval: 200
      title: "Bitswap Stat"
    bwstat:
      enabled: true
      position:
        top: 0
        left: 2
        height: 1
        width: 2
      refreshInterval: 2
      title: "Bandwidth Stat"
`
