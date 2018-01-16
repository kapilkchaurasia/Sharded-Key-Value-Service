package worker

import (
	"net/rpc"
	"DistributedDatabase/src/common"
	"log"
)

type WorkerCliOperation struct {
	client *rpc.Client
}

func StartWorkerCli(masterSrvAddr string,workerSrvAddrs []string) {
	for _, workerSrvAddr := range workerSrvAddrs {
		conn := common.StartRpcCli(masterSrvAddr, "workerCli") //master server location
		workerCliOperation := &WorkerCliOperation{client: conn}
		workerCliOperation.Registry(workerSrvAddr) //worker server location
	}
}

func (t *WorkerCliOperation) Registry(address string){
	var reply string
	args := &common.Args{ address}
	if err := t.client.Call("MasterOperation.Registry", args, &reply); err != nil{
		log.Fatal("error:", err)
	}
}
