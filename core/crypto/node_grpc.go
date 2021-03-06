/*
Copyright IBM Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package crypto

import (
	"crypto/tls"
	"crypto/x509"
	"errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/hyperledger/fabric/core/comm"
)

func (node *nodeImpl) initTLS() error {
	node.debug("Initiliazing TLS...")

	if node.conf.isTLSEnabled() {
		pem, err := node.ks.loadExternalCert(node.conf.getTLSCACertsExternalPath())
		if err != nil {
			node.error("Failed loading TLSCA certificates chain [%s].", err.Error())

			return err
		}

		node.tlsCertPool = x509.NewCertPool()
		ok := node.tlsCertPool.AppendCertsFromPEM(pem)
		if !ok {
			node.error("Failed appending TLSCA certificates chain.")

			return errors.New("Failed appending TLSCA certificates chain.")
		}
		node.debug("Initiliazing TLS...Done")
	} else {
		node.debug("Initiliazing TLS...Disabled!!!")
	}

	return nil
}

func (node *nodeImpl) getClientConn(address string, serverName string) (*grpc.ClientConn, error) {
	node.debug("Dial to addr:[%s], with serverName:[%s]...", address, serverName)

	if node.conf.isTLSEnabled() {
		node.debug("TLS enabled...")

		config := tls.Config{
			InsecureSkipVerify: false,
			RootCAs:            node.tlsCertPool,
			ServerName:         serverName,
		}
		if node.conf.isTLSClientAuthEnabled() {

		}

		return comm.NewClientConnectionWithAddress(address, false, true, credentials.NewTLS(&config))
	}
	node.debug("TLS disabled...")
	return comm.NewClientConnectionWithAddress(address, false, false, nil)
}
