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

package common

// Endpoint represents an endpoint in the communication.
type Endpoint struct {
	IP      string
	Port    uint16
	Mac     string
	Name    string
	Cmdline string
	Proc    string
}

// MakeEndpointPair returns source and destination endpoints from a TCP or IP tuple
// and a command-line tuple.
func MakeEndpointPair(tuple BaseTuple, cmdlineTuple *CmdlineTuple) (src Endpoint, dst Endpoint) {
	src = Endpoint{
		IP:      tuple.SrcIP.String(),
		Port:    tuple.SrcPort,
		Mac:     tuple.SrcMac,
		Proc:    string(cmdlineTuple.Src),
		Cmdline: string(cmdlineTuple.SrcCommand),
	}
	dst = Endpoint{
		IP:      tuple.DstIP.String(),
		Port:    tuple.DstPort,
		Mac:     tuple.DstMac,
		Proc:    string(cmdlineTuple.Dst),
		Cmdline: string(cmdlineTuple.DstCommand),
	}
	return src, dst
}
