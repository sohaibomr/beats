// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package flows

import (
	"sync/atomic"
	"time"
)

type tlsVals struct {
	serverName string
	count      int
}

type biFlow struct {
	id       rawFlowID
	killed   uint32
	createTS time.Time
	ts       time.Time

	dir        flowDirection
	stats      [2]*flowStats
	SYN        int
	tcpopt     TCPOptions
	serverName *tlsVals
	prev, next *biFlow
}

type TCPOptions map[uint32]uint32
type Flow struct {
	TCPOpt     TCPOptions
	stats      *flowStats
	ServerName *tlsVals
}

func newBiFlow(id rawFlowID, ts time.Time, dir flowDirection) *biFlow {
	return &biFlow{
		id:       id,
		ts:       ts,
		createTS: ts,
		dir:      dir,
	}
}

func (f *biFlow) kill() {
	atomic.StoreUint32(&f.killed, 1)
}

func (f *biFlow) isAlive() bool {
	return atomic.LoadUint32(&f.killed) == 0
}

func (f *Flow) AddName(name string) {
	f.ServerName.serverName = name
	f.ServerName.count += 1
	// fmt.Println("Adding server name in flow:", name)
}
