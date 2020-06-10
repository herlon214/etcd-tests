package main

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/herlon214/etcd-tests/pkg/etcd"
	"go.etcd.io/etcd/clientv3"
)

func main() {
	// Env
	caCert := os.Getenv("ETCDCTL_CACERT")
	clientCert := os.Getenv("ETCDCTL_CERT")
	clientKey := os.Getenv("ETCDCTL_KEY")
	endpoints := strings.Split(os.Getenv("ETCDCTL_ENDPOINTS"), ",")

	// Client
	cli, err := etcd.New(endpoints, caCert, clientCert, clientKey)
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	// Create the lease
	res, err := cli.Grant(context.Background(), 5)
	if err != nil {
		log.Fatal(err)
	}

	_, err = cli.Put(context.Background(), "/services/a", "192.168.1.1", clientv3.WithLease(res.ID))
	if err != nil {
		log.Fatal(err)
	}
}
