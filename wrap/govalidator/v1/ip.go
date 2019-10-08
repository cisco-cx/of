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
// The MIT License (MIT)
//
// Copyright (c) 2014 Alex Saskevich
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
	"fmt"

	"github.com/asaskevich/govalidator"

	of "github.com/cisco-cx/of/lib/v1"
)

// IP represents an IPv4 or IPv6 address.
type IP struct {
	ofIP of.IP
}

// NewIP returns a new instance of IP.
func NewIP(ip string) (IP, error) {
	return IP{
		ofIP: of.IP(string(ip)),
	}, nil
}

// IsIPv4 validate that a given IP is an IPv4 address.
func (ip IP) IsIPv4() bool {
	return govalidator.IsIPv4(string(ip.ofIP))
}

// IsIPv6 validate that a given IP is an IPv6 address.
func (ip IP) IsIPv6() bool {
	return govalidator.IsIPv6(string(ip.ofIP))
}

// String implements the fmt.Stringer interface.
func (ip IP) String() string {
	return fmt.Sprintf("%v", ip.ofIP)
}
