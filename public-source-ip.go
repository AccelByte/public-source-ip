// Copyright 2018 AccelByte Inc

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at

// http://www.apache.org/licenses/LICENSE-2.0

// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package public_source_ip

import (
	"net"
	"strings"
)

type PublicSourceIP struct {
	IpNets    []*net.IPNet
	XFFHeader string
}

// New public ip instance
func New(xffHeader string) *PublicSourceIP {
	networks := []*net.IPNet{
		&net.IPNet{net.IPv4(10, 0, 0, 0), net.CIDRMask(8, 32)},
		&net.IPNet{net.IPv4(192, 168, 0, 0), net.CIDRMask(16, 32)},
		&net.IPNet{net.IPv4(172, 16, 0, 0), net.CIDRMask(12, 32)},
		&net.IPNet{net.IPv4(169, 254, 0, 0), net.CIDRMask(16, 32)},
		&net.IPNet{net.IPv4(127, 0, 0, 1), net.CIDRMask(8, 32)},
		// Unique local address
		&net.IPNet{net.IP{0xfc, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, net.CIDRMask(7, 128)},
		// Local addresses
		&net.IPNet{net.IP{0xfe, 0x80, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, net.CIDRMask(10, 128)},
		&net.IPNet{net.IPv6loopback, net.CIDRMask(10, 128)},
	}

	return &PublicSourceIP{
		IpNets:    networks,
		XFFHeader: xffHeader,
	}
}

// isRouted test if an IP is routed on the internet, if it's not then it's a private network
func (publicSourceIP *PublicSourceIP) isRouted(ip net.IP) bool {
	if len(ip) == 0 || ip.Equal(net.IPv4zero) || ip.Equal(net.IPv6zero) {
		return false
	}
	for i := range publicSourceIP.IpNets {
		if publicSourceIP.IpNets[i].Contains(ip) == true {
			return false
		}
	}
	return true
}

// PublicIP determines the public IP used to connect to the first ingress load balancer
func (publicSourceIP *PublicSourceIP) PublicIP() string {
	// Example header:
	// X-Forwarded-For: 10.1.1.2, 203.0.113.195, 70.41.3.18, 150.172.238.178, 192.168.1.1
	nospaces := strings.Replace(publicSourceIP.XFFHeader, " ", "", -1)
	forwards := strings.Split(nospaces, ",")
	for i, _ := range forwards {
		i = len(forwards) - i - 1 // Trace backwards to first external IP address
		ip := net.ParseIP(forwards[i])
		if publicSourceIP.isRouted(ip) == true {
			return ip.String()
		}
	}
	return ""
}
