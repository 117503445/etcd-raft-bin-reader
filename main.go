package main

import (
	"bytes"
	"encoding/gob"
	"flag"
	"fmt"

	"github.com/coreos/etcd/raft/raftpb"
	"github.com/coreos/etcd/wal/walpb"
	"go.etcd.io/etcd/snap"
	"go.etcd.io/etcd/wal"
)

type kv struct {
	Key string
	Val string
}

var (
	walDir  string
	snapDir string
)

func init() {
	flag.StringVar(&walDir, `wal`, `./data/raftexample-1`, `Directory where WAL exists`)
	flag.StringVar(&snapDir, `snapshot`, `./data/raftexample-1-snap`, `Directory where snapshot exists`)
	flag.Parse()
}

func main() {
	fmt.Printf("WAL Directory: %s, Snapshot Dictory: %s\n", walDir, snapDir)
	snapshot, err := loadSnapshot(snapDir)
	if err != nil {
		fmt.Printf("load snapshot failed, err = %v\n", err)
	} else {
		fmt.Println(`========= snapshot =========`)

		fmt.Println(`--------- Metadata ---------`)
		fmt.Printf("Index: %d\n", snapshot.Metadata.Index)
		fmt.Printf("Term: %d\n", snapshot.Metadata.Term)

		fmt.Printf("Conf.nodes: %v\n", snapshot.Metadata.ConfState.Nodes)

		fmt.Println(``)
		fmt.Println(`--------- Data ---------`)
		fmt.Println(string(snapshot.Data))
	}
	fmt.Println(``)

	w, err := openWAL(walDir, snapshot)
	if err != nil {
		fmt.Printf("open wal failed, err = %v\n", err)
	} else {
		fmt.Println(`========= WAL =========`)

		metadata, state, ents, err := w.ReadAll()
		if err != nil {
			panic(err)
		}

		fmt.Println(``)
		fmt.Println(`--------- Metadata ---------`)
		fmt.Printf("metadata: %s\n", metadata)

		fmt.Println(``)
		fmt.Println(`--------- HardState ---------`)
		fmt.Printf("Term: %d\n", state.Term)
		fmt.Printf("Commit: %d\n", state.Commit)
		fmt.Printf("Vote: %d\n", state.Vote)

		fmt.Println(``)
		fmt.Println(`--------- Entries ---------`)
		for i, ent := range ents {
			fmt.Println(``)
			fmt.Printf("Index %d\n", i)
			fmt.Printf("Entry.Index: %d\n", ent.Index)
			fmt.Printf("Entry.Term: %d\n", ent.Term)
			fmt.Printf("Entry.Type: %s\n", ent.Type.String())
			// fmt.Printf("Entry.Data: %s\n", ent.Data)
			var kvData kv
			if err := gob.NewDecoder(bytes.NewBuffer(ent.Data)).Decode(&kvData); err != nil {
				fmt.Println(`ent.Data is not invalid...?`)
			} else {
				fmt.Printf("Entry.Data (gob decoded): %#v\n", kvData)
			}
		}
	}

	fmt.Println(`Done!`)
}

func loadSnapshot(snapDir string) (*raftpb.Snapshot, error) {
	snapshotter := snap.New(snapDir)
	snapshot, err := snapshotter.Load()
	if err != nil {
		return nil, err
	}
	return snapshot, nil
}

func openWAL(walDir string, snapshot *raftpb.Snapshot) (*wal.WAL, error) {
	walsnap := walpb.Snapshot{}
	if snapshot != nil {
		walsnap.Index, walsnap.Term = snapshot.Metadata.Index, snapshot.Metadata.Term
	}
	fmt.Printf("loading WAL at term %d and index %d\n", walsnap.Term, walsnap.Index)
	w, err := wal.Open(walDir, walsnap)
	if err != nil {
		return nil, err
	}

	return w, nil
}
