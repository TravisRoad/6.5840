package kvsrv

import (
	"log"
	"sync"
)

const Debug = false

func DPrintf(format string, a ...interface{}) (n int, err error) {
	if Debug {
		log.Printf(format, a...)
	}
	return
}

type KVServer struct {
	mu sync.Mutex

	// Your definitions here.
	store map[string]string

	IDSet map[int64]ClientValue
}

type ClientValue struct {
	ID    int64
	Value string
}

func (kv *KVServer) Get(args *GetArgs, reply *GetReply) {
	// Your code here.
	kv.mu.Lock()
	defer kv.mu.Unlock()
	// if _, ok := kv.IDSet[args.ID]; ok {
	// 	return
	// }
	// kv.IDSet[args.ID] = struct{}{}
	reply.Value = kv.store[args.Key]
}

func (kv *KVServer) Put(args *PutAppendArgs, reply *PutAppendReply) {
	// Your code here.
	kv.mu.Lock()
	defer kv.mu.Unlock()
	kv.store[args.Key] = args.Value
	reply.Value = args.Value
}

func (kv *KVServer) Append(args *PutAppendArgs, reply *PutAppendReply) {
	// Your code here.
	kv.mu.Lock()
	defer kv.mu.Unlock()
	if v, ok := kv.IDSet[args.ClientID]; ok {
		if v.ID == args.ID {
			reply.Value = v.Value
			return
		}
		// 否则是一个新的请求
	}

	if _, ok := kv.store[args.Key]; !ok {
		kv.store[args.Key] = args.Value
		return
	}
	s := kv.store[args.Key]
	reply.Value = s
	kv.store[args.Key] = s + args.Value
	kv.IDSet[args.ClientID] = ClientValue{
		ID:    args.ID,
		Value: s,
	}
}

func StartKVServer() *KVServer {
	kv := new(KVServer)

	// You may need initialization code here.
	kv.store = make(map[string]string)
	kv.IDSet = make(map[int64]ClientValue)

	return kv
}
