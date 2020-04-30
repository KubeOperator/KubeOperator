/*
 * Tencent is pleased to support the open source community by making TKEStack
 * available.
 *
 * Copyright (C) 2012-2019 Tencent. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License"); you may not use
 * this file except in compliance with the License. You may obtain a copy of the
 * License at
 *
 * https://opensource.org/licenses/Apache-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OF ANY KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations under the License.
 */

package kubeconfig

import (
	"fmt"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

// CreateBasic creates a basic, general KubeConfig object that then can be extended
func CreateBasic(serverURL, clusterName, userName string, caCert []byte) *clientcmdapi.Config {
	// Use the cluster and the username as the context name
	contextName := fmt.Sprintf("%s@%s", userName, clusterName)

	return &clientcmdapi.Config{
		Clusters: map[string]*clientcmdapi.Cluster{
			clusterName: {
				Server:                   serverURL,
				CertificateAuthorityData: caCert,
			},
		},
		Contexts: map[string]*clientcmdapi.Context{
			contextName: {
				Cluster:  clusterName,
				AuthInfo: userName,
			},
		},
		AuthInfos:      map[string]*clientcmdapi.AuthInfo{},
		CurrentContext: contextName,
	}
}

// CreateWithCerts creates a KubeConfig object with access to the API server with client certificates
func CreateWithCerts(serverURL, clusterName, userName string, caCert []byte, clientKey []byte, clientCert []byte) *clientcmdapi.Config {
	config := CreateBasic(serverURL, clusterName, userName, caCert)
	config.AuthInfos[userName] = &clientcmdapi.AuthInfo{
		ClientKeyData:         clientKey,
		ClientCertificateData: clientCert,
	}
	return config
}

// CreateWithToken creates a KubeConfig object with access to the API server with a token
func CreateWithToken(serverURL, clusterName, userName string, caCert []byte, token string) *clientcmdapi.Config {
	config := CreateBasic(serverURL, clusterName, userName, caCert)
	config.AuthInfos[userName] = &clientcmdapi.AuthInfo{
		Token: token,
	}
	return config
}
