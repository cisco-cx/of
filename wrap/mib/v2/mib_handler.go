// Copyright 2019 Cisco Systems, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package v2

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	of "github.com/cisco-cx/of/pkg/v2"
	cache "github.com/cisco-cx/of/wrap/cache/v2"
)

// Load MIBs from different inputs and write
type MIBHandlerer interface {
	// Load MIBs from a JSON file
	LoadJSONFromFile(path string) error
	// Load MIBs from JSON files placed in a dir
	LoadJSONFromDir(dir string) error
	// Load MIBs from cache
	LoadCacheFromFile(path string) error
	// Save MIBs cached into a file
	WriteCacheToFile(path string) error
}

type MIBHandler struct {
	MapMIB map[string]of.MIB
}

func (mh *MIBHandler) LoadJSONFromFile(path string) error {
	jsonFile, err := os.Open(path)
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	var mibJSON map[string]map[string]interface{}
	if err := json.NewDecoder(jsonFile).Decode(&mibJSON); err != nil {
		return err
	}

	for _, entry := range mibJSON {
		if oid, hasOid := entry["oid"]; hasOid {
			var name, description, units string
			name = entry["name"].(string)
			if description, hasDescription := entry["description"]; hasDescription {
				description = description.(string)
			}
			if units, hasUnits := entry["units"]; hasUnits {
				units = units.(string)
			}

			mh.MapMIB[oid.(string)] = of.MIB{
				Name:        name,
				Description: description,
				Units:       units,
			}
		}
	}
	return nil
}

func (mh *MIBHandler) LoadJSONFromDir(path string) error {
	st, err := os.Stat(path)
	if err != nil {
		return err
	}
	if st.IsDir() == false {
		return of.Error(fmt.Sprintf("Path %s is not a directory", path))
	}

	err = filepath.Walk(path, func(filePath string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		return mh.LoadJSONFromFile(filePath)
	})
	return err
}

func (mh *MIBHandler) WriteCacheToFile(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	var cache cache.Cacher
	return cache.Write(f, &mh.MapMIB)
}

func (mh *MIBHandler) LoadCacheFromFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()

	var cache cache.Cacher
	return cache.Read(f, &mh.MapMIB)
}
