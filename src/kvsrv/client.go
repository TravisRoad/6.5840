package kvsrv

import (
	"crypto/rand"
	"math/big"
	"sync"

	"6.5840/labrpc"
)

var (
	IDGen     *IDGenerator
	IDGenOnce sync.Once
)

type IDGenerator struct {
	mu      sync.Mutex
	Current int64
}

func GetIDGenerator() *IDGenerator {
	IDGenOnce.Do(func() {
		IDGen = new(IDGenerator)
	})
	return IDGen
}

func (idg *IDGenerator) NewID() int64 {
	idg.mu.Lock()
	defer idg.mu.Unlock()
	idg.Current++
	return idg.Current
}

func NewID() int64 {
	idg := GetIDGenerator()
	return idg.NewID()
}

type Clerk struct {
	server *labrpc.ClientEnd
	// You will have to modify this struct.
	ID int64
}

func nrand() int64 {
	max := big.NewInt(int64(1) << 62)
	bigx, _ := rand.Int(rand.Reader, max)
	x := bigx.Int64()
	return x
}

func MakeClerk(server *labrpc.ClientEnd) *Clerk {
	ck := new(Clerk)
	ck.server = server
	// You'll have to add code here.
	ck.ID = NewID()
	return ck
}

// fetch the current value for a key.
// returns "" if the key does not exist.
// keeps trying forever in the face of all other errors.
//
// you can send an RPC with code like this:
// ok := ck.server.Call("KVServer.Get", &args, &reply)
//
// the types of args and reply (including whether they are pointers)
// must match the declared types of the RPC handler function's
// arguments. and reply must be passed as a pointer.
func (ck *Clerk) Get(key string) string {
	reply := GetReply{}
	ok := false

	id := NewID()

	for !ok {
		ok = ck.server.Call("KVServer.Get", &GetArgs{ID: id, Key: key}, &reply)
	}

	// You will have to modify this function.
	return reply.Value
}

// shared by Put and Append.
//
// you can send an RPC with code like this:
// ok := ck.server.Call("KVServer."+op, &args, &reply)
//
// the types of args and reply (including whether they are pointers)
// must match the declared types of the RPC handler function's
// arguments. and reply must be passed as a pointer.
func (ck *Clerk) PutAppend(key string, value string, op string) string {
	reply := PutAppendReply{}

	ok := false
	id := NewID()

	for !ok {
		ok = ck.server.Call(
			"KVServer."+op,
			&PutAppendArgs{
				ID:       id,
				ClientID: ck.ID,
				Key:      key,
				Value:    value,
			},
			&reply,
		)
	}

	// You will have to modify this function.
	return reply.Value
}

func (ck *Clerk) Put(key string, value string) {
	ck.PutAppend(key, value, "Put")
}

// Append value to key's value and return that value
func (ck *Clerk) Append(key string, value string) string {
	return ck.PutAppend(key, value, "Append")
}
