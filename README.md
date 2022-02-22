# etcd-raft-bin-reader

Simple Reader of [Etcd Raft](https://github.com/etcd-io/etcd/tree/main/contrib/raftexample) WAL & Snap.

## Getting Started

### Installation

TODO

```sh
go install github.com/117503445/etcd-raft-bin-reader@latest # build from source
```

### Usage

```sh
./etcd-raft-bin-reader -wal data/raftexample-1 -snapshot data/raftexample-1-snap # TODO

go run main.go # dev
```

## License

Distributed under the Apache License. See `LICENSE.txt` for more information.

## Acknowledgments

<https://github.com/nrnrk/raft-wal-viewer> provided the initial version of the code.
