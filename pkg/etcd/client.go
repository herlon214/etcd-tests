package etcd

import (
	"time"

	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/pkg/transport"
)

// New creates a new etcd backend instance
func New(endpoints []string, caPath string, certPath string, certKeyPath string) (*clientv3.Client, error){
	var client *clientv3.Client
	var err error

		tlsInfo := transport.TLSInfo{
			CertFile:      certPath,
			KeyFile:       certKeyPath,
			TrustedCAFile: caPath,
		}
		tlsConfig, err := tlsInfo.ClientConfig()
		if err != nil {
			return nil, err
		}

		client, err = clientv3.New(clientv3.Config{
			Endpoints:   endpoints,
			DialTimeout: 5 * time.Second,
			TLS:         tlsConfig,
		})
		if err != nil {
			return nil, err
		}


	return client, nil
}
