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
// Copyright (c) 2009 The Go Authors. All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are
// met:
//
//    * Redistributions of source code must retain the above copyright
// notice, this list of conditions and the following disclaimer.
//    * Redistributions in binary form must reproduce the above
// copyright notice, this list of conditions and the following disclaimer
// in the documentation and/or other materials provided with the
// distribution.
//    * Neither the name of Google Inc. nor the names of its
// contributors may be used to endorse or promote products derived from
// this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

package v1alpha1

import (
	"fmt"
	"net"

	of "github.com/cisco-cx/of/lib/v1alpha1"
	govalidator "github.com/cisco-cx/of/wrap/govalidator/v1alpha1"
)

// Hostname represents the fully-qualified domain name of a host.
type Hostname struct {
	ofHostname of.Hostname
}

// NewHostname returns a new instance of Hostname.
func NewHostname(h string) (Hostname, error) {
	return Hostname{
		ofHostname: of.Hostname(string(h)),
	}, nil
}

// IPv4() determines the IPv4 addresses associated with a Hostname or
// similar entity.
//
func (h Hostname) IPv4() ([]of.IP, error) {
	results, err := net.LookupHost(string(h.ofHostname))
	if err != nil {
		return []of.IP{}, err
	}
	ipv4s := []of.IP{}
	for _, v := range results {
		ip, err := govalidator.NewIP(v)
		if err != nil {
			return []of.IP{}, err
		}
		if !ip.IsIPv4() {
			continue
		}
		ipv4s = append(ipv4s, of.IP(ip.String()))
	}
	return ipv4s, nil
}

// IPv6() determines the IPv6 addresses associated with a Hostname or
// similar entity.
//
func (h Hostname) IPv6() ([]of.IP, error) {
	results, err := net.LookupHost(string(h.ofHostname))
	if err != nil {
		return []of.IP{}, err
	}
	ipv6s := []of.IP{}
	for _, v := range results {
		ip, err := govalidator.NewIP(v)
		if err != nil {
			return []of.IP{}, err
		}
		if !ip.IsIPv6() {
			continue
		}
		ipv6s = append(ipv6s, of.IP(ip.String()))
	}
	return ipv6s, nil
}

// String implements the fmt.Stringer interface.
func (h Hostname) String() string {
	return fmt.Sprintf("%v", h.ofHostname)
}
