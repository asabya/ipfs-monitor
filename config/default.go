package config

const defaultConfigFile = `tapp:
  colors:
    border:
      focusable: darkslateblue
      focused: orange
      normal: gray
    background: black
    text: white
  grid:
    columns: [32, 32, 32, 32, 90]
    rows: [10, 10, 10, 4, 4, 90]
    background: black
    border: false
  refreshInterval: 1
  widgets:
    peerinfo:
      enabled: true
      position:
        top: 0
        left: 0
        height: 4
        width: 1
      refreshInterval: 0
      title: "Peer Info"
    bootstraplist:
      enabled: true
      position:
        top: 0
        left: 1
        height: 4
        width: 1
      refreshInterval: 0
      title: "Bootstraps"
    repostat:
      enabled: true
      position:
        top: 0
        left: 2
        height: 4
        width: 1
      refreshInterval: 0
      title: "Repo Stats"
`
