package worker

import (
	"fmt"
	"DistributedDatabase/src/common"
	"github.com/peterbourgon/raft"
)

type WorkerSrvOperation struct{
	s *raft.Server
}
var raftRegistry = make(map[int]*raft.Server,0)

func StartWorkerSrv(port int, raftId int ) string {
	op := new(WorkerSrvOperation)
	common.StartRpcServer(op,port,nil,"worker-server") //self location :9091 ||9092 || 9093
	s := RaftServer(uint64(raftId), port)
	raftRegistry[raftId-1]= s
	return  fmt.Sprintf( "%d",port)
}

func (m *WorkerSrvOperation) SetConfig(arg *common.SetConfigArg, reply *string) error {
	SetConfig(raftRegistry[arg.Idx],arg.PeerTuple)
	return nil
}

func (m *WorkerSrvOperation) GetValue(arg *common.GetValueArg, reply *[]byte) error {
	response := make(chan []byte)
	if err := raftRegistry[arg.Idx].Command([]byte(fmt.Sprintf("GET:%s %s",arg)), response); err != nil {
		panic(err) // command not accepted
	}
	*reply = <-response // read value
	return nil
}

func (m *WorkerSrvOperation) PutValue(arg *common.PutValueArg, reply *string) error {
	response := make(chan []byte)
	if err := raftRegistry[arg.Idx].Command([]byte(fmt.Sprintf("PUT:%s %s",arg.Key,arg.Value)), response); err != nil {
		panic(err) // command not accepted
	}
	<-response // read value
	return nil
}