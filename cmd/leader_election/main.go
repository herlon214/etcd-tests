package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/herlon214/etcd-tests/pkg/etcd"
	"go.etcd.io/etcd/clientv3/concurrency"
)

func main() {
	var name = flag.String("name", "", "node name")
	flag.Parse()

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

	// Create a new session
	sess, err := concurrency.NewSession(cli)
	if err != nil {
		log.Fatal(err)
	}

	for {
		// Create a new election
		e := concurrency.NewElection(sess, "/master-node")

		// Elect a leader
		if err := e.Campaign(context.Background(), *name); err != nil {
			log.Fatal(err)
		}

		fmt.Println("Doing something...", *name)
		time.Sleep(5 * time.Second)
		if err := e.Resign(context.Background()); err != nil {
			log.Fatal(err)
		}

		fmt.Println("Resigned ", *name)
	}

}
