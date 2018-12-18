package balance

import (
	"encoding/json"
	etcd3 "github.com/coreos/etcd/clientv3"
	"golang.org/x/net/context"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/naming"
	log "github.com/sirupsen/logrus"
	"strings"
)

// EtcdWatcher is the implementation of grpc.naming.Watcher
type EtcdWatcher struct {
	key     string
	client  *etcd3.Client
	updates []*naming.Update
	ctx     context.Context
	cancel  context.CancelFunc
}

func (w *EtcdWatcher) Close() {
	w.cancel()
}

func newEtcdWatcher(key string, cli *etcd3.Client) naming.Watcher {
	ctx, cancel := context.WithCancel(context.Background())
	w := &EtcdWatcher{
		key:     key,
		client:  cli,
		ctx:     ctx,
		updates: make([]*naming.Update, 0),
		cancel:  cancel,
	}
	return w
}

func (w *EtcdWatcher) Next() ([]*naming.Update, error) {
	updates := make([]*naming.Update, 0)

	if len(w.updates) == 0 {
		// query addresses from etcd
		resp, err := w.client.Get(w.ctx, w.key, etcd3.WithPrefix())
		if err == nil {
			addrs := extractAddrs(resp)
			if len(addrs) > 0 {
				for _, v := range addrs {
					updates = append(updates, &naming.Update{Op: naming.Add, Addr: v.Addr, Metadata: &v.Metadata})
				}
				w.updates = updates
				return updates, nil
			}
		} else {
			grpclog.Println("Etcd Watcher Get key error:", err)
		}
	}

	// generate etcd Watcher
	rch := w.client.Watch(w.ctx, w.key, etcd3.WithPrefix())
	select {
	case wresp := <-rch:
		for _, ev := range wresp.Events {
			switch ev.Type {

			case etcd3.EventTypePut:
				nodeData := NodeData{}
				err := json.Unmarshal([]byte(ev.Kv.Value), &nodeData)
				if err != nil {
					log.Debugf("Parse node data error:", err)
					continue
				}
				updates = append(updates, &naming.Update{Op: naming.Add, Addr: nodeData.Addr, Metadata: &nodeData.Metadata})
			case etcd3.EventTypeDelete:
				key := string(ev.Kv.Key)
				key = key[strings.LastIndex(key, "/")+1:]
				updates = append(updates, &naming.Update{Op: naming.Delete, Addr: key, Metadata: nil})
			}
		}
	}
	return updates, nil
}

func extractAddrs(resp *etcd3.GetResponse) []NodeData {
	addrs := []NodeData{}

	if resp == nil || resp.Kvs == nil {
		return addrs
	}

	for i := range resp.Kvs {
		if v := resp.Kvs[i].Value; v != nil {
			nodeData := NodeData{}
			err := json.Unmarshal(v, &nodeData)
			if err != nil {
				grpclog.Println("Parse node data error:", err)
				continue
			}
			addrs = append(addrs, nodeData)
		}
	}

	return addrs
}