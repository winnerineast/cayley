// Copyright 2016 The Cayley Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package bolt

import (
	"os"
	"path/filepath"

	"github.com/cayleygraph/cayley/clog"
	"github.com/cayleygraph/cayley/graph"
	hkv "github.com/nwca/hidalgo/kv"
	"github.com/nwca/hidalgo/kv/bolt"
)

const (
	Type = bolt.Name
)

func getBoltFile(cfgpath string) string {
	return filepath.Join(cfgpath, "indexes.bolt")
}

func Create(path string, _ graph.Options) (hkv.KV, error) {
	err := os.MkdirAll(path, 0700)
	if err != nil {
		return nil, err
	}
	db, err := bolt.Open(getBoltFile(path), nil)
	if err != nil {
		clog.Errorf("Error: couldn't create Bolt database: %v", err)
		return nil, err
	}
	return db, nil
}

func Open(path string, opt graph.Options) (hkv.KV, error) {
	db, err := bolt.Open(getBoltFile(path), nil)
	if err != nil {
		clog.Errorf("Error, couldn't open! %v", err)
		return nil, err
	}
	bdb := db.DB()
	// BoolKey returns false on non-existence. IE, Sync by default.
	bdb.NoSync, err = opt.BoolKey("nosync", false)
	if err != nil {
		db.Close()
		return nil, err
	}
	if bdb.NoSync {
		clog.Infof("Running in nosync mode")
	}
	return db, nil
}
