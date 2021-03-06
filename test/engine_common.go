// Copyright 2016 Keybase Inc. All rights reserved.
// Use of this source code is governed by a BSD
// license that can be found in the LICENSE file.

package test

import (
	"sort"
	"testing"

	"github.com/keybase/client/go/libkb"
	"github.com/keybase/kbfs/libkbfs"
)

func setBlockSizes(t testing.TB, config libkbfs.Config, blockSize, blockChangeSize int64) {
	// Set the block sizes, if any
	if blockSize > 0 || blockChangeSize > 0 {
		if blockSize == 0 {
			blockSize = 512 * 1024
		}
		if blockChangeSize < 0 {
			t.Fatal("Can't handle negative blockChangeSize")
		}
		if blockChangeSize == 0 {
			blockChangeSize = 8 * 1024
		}
		bsplit, err := libkbfs.NewBlockSplitterSimple(blockSize,
			uint64(blockChangeSize), config.Codec())
		if err != nil {
			t.Fatalf("Couldn't make block splitter for block size %d,"+
				" blockChangeSize %d: %v", blockSize, blockChangeSize, err)
		}
		config.SetBlockSplitter(bsplit)
	}
}

func maybeSetBw(t testing.TB, config libkbfs.Config, bwKBps int) {
	if bwKBps > 0 {
		config.SetBlockOps(libkbfs.NewBlockOpsConstrained(
			config.BlockOps(), bwKBps))
		// Looks like we're testing big transfers, so let's do
		// background flushes.
		config.SetDoBackgroundFlushes(true)
	}
}

func makeTeams(t testing.TB, config libkbfs.Config, e Engine, teams teamMap,
	users map[libkb.NormalizedUsername]User) {
	teamNames := make([]libkb.NormalizedUsername, 0, len(teams))
	for name := range teams {
		teamNames = append(teamNames, name)
	}
	sort.Slice(teamNames, func(i, j int) bool {
		return string(teamNames[i]) < string(teamNames[j])
	}) // make sure root names go first.
	infos := libkbfs.AddEmptyTeamsForTestOrBust(t, config, teamNames...)
	for i, name := range teamNames {
		members := teams[name]
		tid := infos[i].TID
		for _, w := range members.writers {
			libkbfs.AddTeamWriterForTestOrBust(t, config, tid,
				e.GetUID(users[w]))
		}
		for _, r := range members.readers {
			libkbfs.AddTeamReaderForTestOrBust(t, config, tid,
				e.GetUID(users[r]))
		}
	}
}
