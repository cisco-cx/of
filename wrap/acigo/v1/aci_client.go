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
//
// This work incorporates works covered by the following notices:
//
// MIT License
//
// Copyright (c) 2016 udhos
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package v1

import (
	"strings"

	"github.com/udhos/acigo/aci"

	of "github.com/cisco-cx/of/pkg/v1"
	logger "github.com/cisco-cx/of/wrap/logrus/v1"
	mapstructure "github.com/cisco-cx/of/wrap/mapstructure/v1"
)

// ACIClient represents an instance of the acigo ACI API client.
//
// ACIClient implements the of.ACIClient interface.
type ACIClient struct {
	client *aci.Client
	Log    *logger.Logger
}

// NewACIClient returns a new instance of ACIClient configured by an
// of.ACIClientConfig struct.
func NewACIClient(cfg of.ACIClientConfig, log *logger.Logger) (*ACIClient, error) {
	// Convert of.ACIClientConfig to aci.ClientOptions.
	opts := aci.ClientOptions{
		Hosts: cfg.Hosts,
		User:  cfg.User,
		Pass:  cfg.Pass,
		Debug: cfg.Debug,
	}
	// Configure the new internal client.
	client, err := aci.New(opts)
	if err != nil {
		return &ACIClient{}, err
	}
	return &ACIClient{client: client, Log: log}, nil
}

// Login opens a new API session.
func (c *ACIClient) Login() error {
	return c.client.Login()
}

// Faults returns the API's "fault table" as a slice of of.Map objects.
func (c *ACIClient) Faults() ([]of.Map, error) {
	list, err := c.client.FaultList()
	if err != nil {
		return []of.Map{}, err
	}
	mm := make([]of.Map, len(list))
	for i, v := range list {
		c.Log.Tracef("Recieved fault from ACI: %+v\n", v)
		mapstructure.NewMap(v).DecodeMap(&mm[i])
	}
	return mm, nil
}

// NodeList retrieves the list of top level system elements (APICs, spines, leaves).
func (c *ACIClient) NodeList() (map[string]map[string]interface{}, error) {
	nodeMap := make(map[string]map[string]interface{})
	nodes, err := c.client.NodeList()
	if err != nil {
		return nodeMap, err
	}
	for _, node := range nodes {
		c.Log.Tracef("Recieved node from ACI: %+v\n", node)
		dn := strings.Replace(node["dn"].(string), "/sys", "", -1)
		nodeMap[dn] = node
	}
	return nodeMap, nil
}

// Logout closes the API session.
func (c *ACIClient) Logout() error {
	return c.client.Logout()
}

// Allow set the ACI host on runtime
func (c *ACIClient) SetHost(host string) {
	c.client.Opt.Hosts = []string{host}
}

// Get the current ACI host
func (c *ACIClient) GetHost() string {
	if c != nil {
		if len(c.client.Opt.Hosts) >= 1 {
			return c.client.Opt.Hosts[0]
		}
	}
	return ""
}
