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
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsNotRoutedWithOwnInternalNetworkList(t *testing.T) {
	for _, network := range networks {
		assert.Falsef(t, isRouted(network.IP), "non-routable network %v", network.IP.String())
	}
}

func TestIsNotRoutedWithAddressXFFHeaders(t *testing.T) {
	xff := []struct {
		header string
		realip string
	}{
		{"10.1.2.3, 206.27.34.1, 192.168.2.200", "206.27.34.1"},
		{"172.16.255.255, 10.4.0.0.1, 192.1.4.1, 202.0.0.4, 10.1.1.1", "202.0.0.4"},
		{"192.168.1.100,  172.16.1.4, 54.0.1.53, 10.1.4.3, 192.168.1.1", "54.0.1.53"},
		{"l27.0.0.2", ""},
		{"169.123.23.3, 2001:0db8:85a3:0000:0000:8a2e:0370:7334, 192.168.1.1", "2001:db8:85a3::8a2e:370:7334"},
		{"::", ""},
		{"", ""},
		{"fd00::, 2001:0db8:85a3:0000:0000:8a2e:0370:7334, 192.4.1.56", "192.4.1.56"},
	}
	for i := range xff {
		req, err := http.NewRequest("GET", "/", nil)
		req.Header.Set("X-Forwarded-For", xff[i].header)
		ip := PublicIP(req)

		assert.Nil(t, err, "unable to create new request")
		assert.Truef(t, ip == xff[i].realip, "real ip result expected \"%v\", but got \"%v\"", xff[i].realip, ip)
	}
}

func TestIsNotRoutedWithInvalidAddressXFFHeaders(t *testing.T) {
	xff := []struct {
		header string
		realip string
	}{
		{"10.1.2.3206.27.34.1 * * * 192.168.2.200", ""},
		{"172.16.255.999 10.4.0.0.1 192.1.4.1 202.0.0.4 10.1.1.1", ""},
		{"172.16.255.999 10.4.0.0.1 192.1.4.1, 202.0.0.4, 10.1.1.1", "202.0.0.4"},
		{"****", ""},
	}
	for i := range xff {
		req, err := http.NewRequest("GET", "/", nil)
		req.Header.Set("X-Forwarded-For", xff[i].header)
		ip := PublicIP(req)

		assert.Nil(t, err, "unable to create new request")
		assert.Truef(t, ip == xff[i].realip, "real ip result expected \"%v\", but got \"%v\"", xff[i].realip, ip)
	}
}
