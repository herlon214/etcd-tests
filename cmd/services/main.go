package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/herlon214/etcd-tests/pkg/etcd"
	"go.etcd.io/etcd/clientv3"
)

func main() {
	// Mutex to change svc list
	var mx sync.Mutex

	// Service list
	svcList := make(map[string]string, 0)

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

	// Show the actual list of services
	go func() {
		for {
			mx.Lock()
			fmt.Println("---------------------------------------")
			for service, ip := range svcList {
				fmt.Println(service, " -> ", ip)
			}
			fmt.Println("---------------------------------------")
			mx.Unlock()

			time.Sleep(time.Second)
		}

	}()

	// Watch keys
	res := cli.Watch(context.Background(), "/services", clientv3.WithPrefix())
	for wresp := range res {
		for _, ev := range wresp.Events {
			mx.Lock()
			if ev.Type == clientv3.EventTypeDelete {
				delete(svcList, string(ev.Kv.Key))
			} else {
				svcList[string(ev.Kv.Key)] = string(ev.Kv.Value)
			}
			mx.Unlock()
		}
	}

}
