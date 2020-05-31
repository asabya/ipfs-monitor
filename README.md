# ipfsmon

![ipfs-monitor](https://user-images.githubusercontent.com/15252513/83339816-b8871c80-a2ee-11ea-92a5-cad5e92d8b08.png)

> A terminal application for [IPFS](https://ipfs.io) http api.
>
> You need the command line to be running an IPFS node (go-ipfs) to monitor different features.

![ipfs-monitor](https://user-images.githubusercontent.com/15252513/83355299-4b21cd00-a37c-11ea-8a02-8e7d7019205b.png)

IPFS Monitor allows you to monitor the behavior of your IPFS Node without having to bother with different commands.

> ⚠ Please note that this version is not stable yet and will change.

**Download the latest release**

- Mac - [ipfs-monitor](https://github.com/Sab94/ipfs-monitor/releases/download/v0.1.1/ipfsmon_darwin_amd64_0.1.1)
- Windows - [ipfs-monitor.exe](https://github.com/Sab94/ipfs-monitor/releases/download/v0.1.1/ipfsmon_windows_amd64_0.1.1.exe)
- Linux - [ipfs-monitor](https://github.com/Sab94/ipfs-monitor/releases/download/v0.1.1/ipfsmon_linux_amd64_0.1.1)

#### Download and Compile

```
$ git clone https://github.com/Sab94/ipfs-monitor.git

$ cd ipfs-monitor
$ make install
```

Alternatively, you can run `make build` to build the ipfsmon binary (storing it in `cmd/`) without installing it.

##### Cross Compiling

Compiling for a different platform is as simple as running:

```
make build GOOS=TargetOS GOARCH=TargetArchitecture
```

#### Running

```
ipfsmon

or

ipfsmon http://localhost:5001
```

It connects to default ipfs api `http://localhost:5001`, if you are running ipfs on a different port
pass the api url as first argument

**TODO list**
- [ ] add more modules (swarm-addrs, bitswap-ledger, etc...)
- [ ] add tests
- [ ] add a guide for adding more modules
- [ ] auto update on `config.yml` change
- [ ] add docs and comments

> **Note : This project is highly inspired by [wtfutil/wtf](https://github.com/wtfutil/wtf)

## Project52

It is one of my [project 52](https://github.com/Sab94/project52).
