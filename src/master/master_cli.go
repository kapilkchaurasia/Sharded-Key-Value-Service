package master

import (
	"net/rpc"
	"DistributedDatabase/src/common"
	"log"
)

type MasterCliOperation struct {
	Client []*rpc.Client
}
var masterCliOperation = MasterCliOperation{Client:nil}

func StartMasterCli() { // location : 9091,9092,9093
	client := make([]*rpc.Client,0)
	for _, addr := range WorkerSrvAddr {
		conn := common.StartRpcCli(addr, "masterCli")
		client = append(client, conn)
	}
	masterCliOperation = MasterCliOperation{Client:client}
}

func setConfig(configs [][]string) {
	for idx, client := range masterCliOperation.Client  {
		if err :=client.Call("WorkerSrvOperation.SetConfig",&common.SetConfigArg{idx,configs[idx]},
		nil);err !=nil {
			log.Printf("error in setConfig for RF Op on worker-instances %s",masterCliOperation.Client[idx])
		}
	}
}

func putValue(key string, value string) {
		if err :=masterCliOperation.Client[hashingFunc(key)].Call("WorkerSrvOperation.PutValueArg",&common.PutValueArg{key,value,hashingFunc(key)},
			nil);err !=nil {
			log.Printf("error in putValue Op on worker-instances %s",masterCliOperation.Client[hashingFunc(key)])
		}
}

func getValue(key string) []byte{
	var reply []byte
	if err :=masterCliOperation.Client[hashingFunc(key)].Call("WorkerSrvOperation.GetValueArg",&common.GetValueArg{key,hashingFunc(key)},
		&reply);err !=nil {
		log.Printf("error in getValue Op on worker-instances %s",masterCliOperation.Client[hashingFunc(key)])
		return nil
	}
	return reply
}