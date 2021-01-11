/*
 * Copyright 2018 AccelByte Inc
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package public_source_ip

import (
	"net"
	"net/http"
	"strings"
)

type ipNets []*net.IPNet

var (
	// Local networks not routed on the internet (these change occasionally, esp with IPv6)
	networks = ipNets{
		&net.IPNet{IP: net.IPv4(10, 0, 0, 0), Mask: net.CIDRMask(8, 32)},
		&net.IPNet{IP: net.IPv4(192, 168, 0, 0), Mask: net.CIDRMask(16, 32)},
		&net.IPNet{IP: net.IPv4(172, 16, 0, 0), Mask: net.CIDRMask(12, 32)},
		&net.IPNet{IP: net.IPv4(169, 254, 0, 0), Mask: net.CIDRMask(16, 32)},
		&net.IPNet{IP: net.IPv4(127, 0, 0, 1), Mask: net.CIDRMask(8, 32)},
		// Unique local address
		&net.IPNet{IP: net.IP{0xfc, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, Mask: net.CIDRMask(7, 128)},
		// Local addresses
		&net.IPNet{IP: net.IP{0xfe, 0x80, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}, Mask: net.CIDRMask(10, 128)},
		&net.IPNet{IP: net.IPv6loopback, Mask: net.CIDRMask(10, 128)},
		// Shared Address Space, ref https://tools.ietf.org/html/rfc6598
		&net.IPNet{IP: net.IPv4(100, 64, 0, 0), Mask: net.CIDRMask(10, 32)},
	}
)

// isRouted test if an IP is routed on the internet, if it's not then it's a private network
func isRouted(ip net.IP) bool {
	if len(ip) == 0 || ip.Equal(net.IPv4zero) || ip.Equal(net.IPv6zero) {
		return false
	}
	for i := range networks {
		if networks[i].Contains(ip) == true {
			return false
		}
	}
	return true
}

// PublicIP determines the public IP used to connect to the first ingress load balancer
func PublicIP(request *http.Request) string {
	// Example header:
	// X-Forwarded-For: 10.1.1.2, 203.0.113.195, 70.41.3.18, 150.172.238.178, 192.168.1.1
	xffHeader := request.Header.Get("X-Forwarded-For")
	nospaces := strings.Replace(xffHeader, " ", "", -1)
	forwards := strings.Split(nospaces, ",")

	for i := range forwards {
		ip := net.ParseIP(forwards[i])
		if isRouted(ip) == true {
			return ip.String()
		}
	}
	return ""
}
