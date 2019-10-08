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

package v1

// IP represents an IPv4 or IPv6 address.
type IP string

// IPv4Finder is an interface that represents the ability to determine the IPv4
// addresses asssociated with a Hostname or similar entity.
type IPv4Finder interface {
	IPv4() ([]IP, error)
}

// IPv4Validator is an interface that represents the ability to validate that
// a given IP is an IPv4 address.
type IPv4Validator interface {
	IsIPv4() bool
}

// IPv6Finder is an interface that represents the ability to determine the IPv6
// addresses asssociated with a Hostname or similar entity.
type IPv6Finder interface {
	IPv6() ([]IP, error)
}

// IPv6Validator is an interface that represents the ability to validate that
// a given IP is an IPv6 address.
type IPv6Validator interface {
	IsIPv6() bool
}
