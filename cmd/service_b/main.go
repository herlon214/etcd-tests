package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

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
	lease, err := cli.Grant(context.Background(), 5)
	if err != nil {
		log.Fatal(err)
	}

	// Set a new key
	_, err = cli.Put(context.Background(), "/services/b", "192.168.1.1", clientv3.WithLease(lease.ID))
	if err != nil {
		log.Fatal(err)
	}

	// Keep the lease alive
	for {
		ka, err := cli.KeepAliveOnce(context.Background(), lease.ID)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(ka.TTL)

		time.Sleep(time.Second * 2)
	}
}
