/*
 * Copyright © 2015-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * @author		Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @copyright 	2015-2018 Aeneas Rekkas <aeneas+oss@aeneas.io>
 * @license 	Apache-2.0
 */
package herodot

// statusCodeCarrier can be implemented by an error to support setting status codes in the error itself.
type statusCodeCarrier interface {
	// StatusCode returns the status code of this error.
	StatusCode() int
}

// requestIDCarrier can be implemented by an error to support error contexts.
type requestIDCarrier interface {
	// RequestID returns the ID of the request that caused the error, if applicable.
	RequestID() string
}

// reasonCarrier can be implemented by an error to support error contexts.
type reasonCarrier interface {
	// Reason returns the reason for the error, if applicable.
	Reason() string
}

// debugCarrier can be implemented by an error to support error contexts.
type debugCarrier interface {
	// Debug returns debugging information for the error, if applicable.
	Debug() string
}

// statusCarrier can be implemented by an error to support error contexts.
type statusCarrier interface {
	// ID returns the error id, if applicable.
	Status() string
}

// detailsCarrier can be implemented by an error to support error contexts.
type detailsCarrier interface {
	// Details returns details on the error, if applicable.
	Details() map[string][]interface{}
}
