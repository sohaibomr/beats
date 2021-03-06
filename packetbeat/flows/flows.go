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
	"fmt"
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/packetbeat/config"
)

type Flows struct {
	worker     *worker
	table      *flowMetaTable
	counterReg *counterReg
}

// Reporter callback type, to report flow events to.
type Reporter func([]beat.Event)

var debugf = logp.MakeDebug("flows")

const (
	defaultTimeout = 30 * time.Second
	defaultPeriod  = 10 * time.Second
)

func NewFlows(pub Reporter, config *config.Flows) (*Flows, error) {
	fmt.Println("In new flows")
	duration := func(s string, d time.Duration) (time.Duration, error) {
		if s == "" {
			return d, nil
		}
		return time.ParseDuration(s)
	}

	timeout, err := duration(config.Timeout, defaultTimeout)
	if err != nil {
		logp.Err("failed to parse flow timeout: %v", err)
		return nil, err
	}

	period, err := duration(config.Period, defaultPeriod)
	if err != nil {
		logp.Err("failed to parse period: %v", err)
		return nil, err
	}

	table := &flowMetaTable{
		table: make(map[flowIDMeta]*flowTable),
	}

	counter := &counterReg{}

	// below returns a struct Worker, which has run func(), which is bassically a periodic function which runs after a
	// period and creates new events from the flowMetaTables
	worker, err := newFlowsWorker(pub, table, counter, timeout, period)
	if err != nil {
		logp.Err("failed to configure flows processing intervals: %v", err)
		return nil, err
	}
	// fmt.Println("flow table", table)

	return &Flows{
		table:      table,
		worker:     worker,
		counterReg: counter,
	}, nil
}

func (f *Flows) Lock() {
	debugf("lock flows")
	f.table.Lock()
}

func (f *Flows) Unlock() {
	debugf("unlock flows")
	f.table.Unlock()
}

func (f *Flows) Get(id *FlowID) *Flow {
	debugf("get flow")
	// fmt.Println(id.flow.stats)
	if id.flow.stats == nil {
		// fmt.Println("In flow.stats")
		debugf("lookup flow: %v => %v", id.flowIDMeta, id.flowID)
		id.flow = f.table.get(id, f.counterReg)
	}
	return &id.flow
}

// for calculating rtt by T.S
func (f *Flows) GetSYN(id *FlowID) int {
	t := f.table.table[id.flowIDMeta]
	bf := t.table[string(id.flowID)]
	return bf.SYN
}

func (f *Flows) AddSYN(id *FlowID) {
	t := f.table.table[id.flowIDMeta]
	t.mutex.Lock()
	defer t.mutex.Unlock()
	bf := t.table[string(id.flowID)]
	bf.SYN = bf.SYN + 1
	return
}
func (f *Flows) RemoveSYN(id *FlowID) {
	t := f.table.table[id.flowIDMeta]
	t.mutex.Lock()
	defer t.mutex.Unlock()
	bf := t.table[string(id.flowID)]
	bf.SYN = 0
	return
}

func (f *Flows) AddTCPOpt(id *FlowID, tsval uint32, tsecr uint32) {

	t := f.table.table[id.flowIDMeta]
	t.mutex.Lock()
	defer t.mutex.Unlock()

	bf := t.table[string(id.flowID)]
	bf.tcpopt[tsval] = tsecr
	return
}

func (f *Flows) DelTCPOpt(id *FlowID, tsecr uint32) {

	t := f.table.table[id.flowIDMeta]
	t.mutex.Lock()
	defer t.mutex.Unlock()

	bf := t.table[string(id.flowID)]
	//fmt.Println("stats", bf.stats)

	delete(bf.tcpopt, tsecr)
	// fmt.Println()
	// fmt.Println("In DelOpt", bf.tcpopt)

	return
}

// ends

// to add server name in biflow
func (f *Flows) AddServerName(id *FlowID, name string) {
	// fmt.Println("In ad server name")
	// t := f.table.table[id.flowIDMeta]
	// t.mutex.Lock()
	// defer t.mutex.Unlock()

	// bf := t.table[string(id.flowID)]
	// fmt.Println("Server name before adding in biflow", name)
	// bf.serverName = name
	fmt.Println(string(id.flowID))
}

func (f *Flows) Start() {
	f.worker.Start()
}

func (f *Flows) Stop() {
	f.worker.Stop()
}

func (f *Flows) NewInt(name string) (*Int, error) {
	return f.counterReg.newInt(name)
}

func (f *Flows) NewUint(name string) (*Uint, error) {
	return f.counterReg.newUint(name)
}

func (f *Flows) NewFloat(name string) (*Float, error) {
	return f.counterReg.newFloat(name)
}
