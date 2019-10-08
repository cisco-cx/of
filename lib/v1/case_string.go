// Copyright 2019 Cisco Systems, Inc.
//
// This work incorporates works covered by the following notices:
//
// The MIT License (MIT)
//
// Copyright (c) 2015 Ian Coleman
// Copyright (c) 2018 Ma_124, <github.com/Ma124>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, Subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or Substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package v1

// CaseString represents a string that will be commonly used in the receiver of
// the CaseConverter interface.
type CaseString string

// CaseConverter is an interface that represents the ability to return
// a case-converted string. For example, "this-case" -> "ThatCase".
//
// CaseConverter is based on the functions implemented in
// `github.com/iancoleman/strcase`.
type CaseConverter interface {
	ToSnake() string
}
