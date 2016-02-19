/*
 *  Licensed to Wikifeat under one or more contributor license agreements.
 *  See the LICENSE.txt file distributed with this work for additional information
 *  regarding copyright ownership.
 *
 *  Redistribution and use in source and binary forms, with or without
 *  modification, are permitted provided that the following conditions are met:
 *
 *  * Redistributions of source code must retain the above copyright notice,
 *  this list of conditions and the following disclaimer.
 *  * Redistributions in binary form must reproduce the above copyright
 *  notice, this list of conditions and the following disclaimer in the
 *  documentation and/or other materials provided with the distribution.
 *  * Neither the name of Wikifeat nor the names of its contributors may be used
 *  to endorse or promote products derived from this software without
 *  specific prior written permission.
 *
 *  THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
 *  AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 *  IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
 *  ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE
 *  LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
 *  CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
 *  SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
 *  INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
 *  CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
 *  ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
 *  POSSIBILITY OF SUCH DAMAGE.
 */
package config

import (
	"flag"
)

type DefaultCmdLine struct {
	HostName         string
	NodeId           string
	Port             string
	UseSSL           bool
	SSLCertFile      string
	SSLKeyFile       string
	RegistryLocation string
}

func ParseCmdParams(defaults DefaultCmdLine) {
	hostName := flag.String("hostName", defaults.HostName, "The host name for this instance")
	nodeId := flag.String("nodeId", defaults.NodeId, "The node Id for this instance")
	port := flag.String("port", defaults.Port, "The port number for this instance")
	useSSL := flag.Bool("useSSL", defaults.UseSSL, "use SSL")
	sslCertFile := flag.String("sslCertFile", defaults.SSLCertFile, "The SSL certificate file")
	sslKeyFile := flag.String("sslKeyFile", defaults.SSLKeyFile, "The SSL key file")
	registryLocation := flag.String("registryLocation", defaults.RegistryLocation, "URL for etcd")
	flag.Parse()
	Service.DomainName = *hostName
	Service.NodeId = *nodeId
	Service.Port = *port
	Service.UseSSL = *useSSL
	Service.SSLCertFile = *sslCertFile
	Service.SSLKeyFile = *sslKeyFile
	Service.RegistryLocation = *registryLocation
}
