package node

import (
	"math/rand"
	"net/rpc"
	"time"
)

type Downloader struct {
	node *Node
}

func (d *Downloader) getFile(file FileInfo) error {
	var randomNode []NodeInfo
	err := d.node.rpcConn.Call("Server.GetNode", MinNumRoute*10, &randomNode)
	if err != nil {
		return err
	}

	numReq := file.Size/FileChunkSize + 1
	reqPerNode := len(file.Nodes) / numReq

	writeFile(file)
	d.node.fileStatus[file.Fname] = make([]byte, numReq)

	var data Data
	var chunkInfo Chunk
	chunkInfo.Length = numReq
	data.Data = chunkInfo
	data.FInfo = file

	Log.Printf("GetFile - FileName - [%s], FileSize - [%d]", file.Fname, file.Size)
	Log.Printf("        - ChunkLength - [%d], NumReqPerNode - [%d]", numReq, reqPerNode)

	index := 0
	for _, node := range file.Nodes {
		for i := 0; i < reqPerNode; i++ {
			Log.Printf("GetFile - Requesting [%d/%d]", index, numReq)
			data.Data.Index = index
			data.PublicKey = d.node.dataPublicKey
			route, operation := d.generatePath(node, randomNode)
			aesKey, encrpytedData, err := EncryptStruct(data, node.PublicKey)
			if err != nil {
				Log.Fatal("GetFile - Encrypting Data Failed")
			}
			dataBox := DataBox{aesKey, encrpytedData}
			message := Message{operation.Next, route.Next, dataBox}
			//message := Message{operation.Next, route.Next, encrpytedData}

			reply := false
			conn, _ := rpc.Dial("tcp", route.Dst)
			conn.Call("Node.Incoming", message, &reply)
			conn.Close()

			index++
		}
	}

	Log.Printf("GetFile - Requesting [%d/%d]", index, numReq)
	node := file.Nodes[0]
	data.Data.Index = index
	data.PublicKey = d.node.dataPublicKey
	route, operation := d.generatePath(node, randomNode)
	aesKey, encrpytedData, err := EncryptStruct(data, node.PublicKey)
	//encrpytedData, err := EncryptStruct(data, node.PublicKey)
	if err != nil {
		Log.Fatal("GetFile - Encrypting Data Failed")
	}
	dataBox := DataBox{aesKey, encrpytedData}
	message := Message{operation.Next, route.Next, dataBox}

	//message := Message{aesKey,operation.Next, route.Next, encrpytedData}

	reply := false
	conn, _ := rpc.Dial("tcp", route.Dst)
	conn.Call("Node.Incoming", message, &reply)
	conn.Close()

	for {
		time.Sleep(10 * time.Second)
		complete := true
		for _, status := range d.node.fileStatus[file.Fname] {
			if status == 0 {
				complete = false
				break
			}
		}

		if complete {
			_, err := doneWriting(file)
			return err
		}
	}
}

func (d *Downloader) generatePath(dst NodeInfo, randomNode []NodeInfo) (Route, Operation) {

	//route := Route{
	//	Dst:  d.node.listener.Addr().String(),
	//	NextMsg: nextMsg,
	//}
	//
	//operation := Operation{
	//	Op:   END,
	//	NextMsg: nextMsg,
	//}

	route := Route{}
	route.Dst = d.node.listener.Addr().String()

	operation := Operation{}
	operation.Op = END

	layerMessage(&route, &operation, randomNode)

	aesKey, next, err := EncryptStruct(route, dst.PublicKey)
	if err != nil {
		Log.Fatal("GetFile - Encrypting Data Failed")
	}

	route.Dst = dst.Addr
	route.Next = RouteBox{aesKey, next}

	//route.Dst = dst.Addr
	//route.Next = next

	aesKey, next, err = EncryptStruct(operation, dst.PublicKey)
	if err != nil {
		Log.Fatal("GetFile - Encrypting Data Failed")
	}

	operation.Op = GETFILE
	operation.Next = OpBox{aesKey, next}

	//operation.Op = GETFILE
	//operation.Next = next

	layerMessage(&route, &operation, randomNode)

	return route, operation
}

func layerMessage(route *Route, operation *Operation, randomNode []NodeInfo) {
	rand.Seed(time.Now().Unix())
	length := len(randomNode)
	for i := 1; i < MinNumRoute; i++ {
		n := rand.Int() % length
		aesKey, encryptedRoute, err := EncryptStruct(*route, randomNode[n].PublicKey)
		if err != nil {
			Log.Fatal("GetFile - Encrypting Data Failed")
		}
		routeBox := RouteBox{aesKey, encryptedRoute}
		route.Next = routeBox

		//route.Next = encryptedRoute
		//route.Dst = randomNode[n].Addr

		aesKey, encryptedOperation, err := EncryptStruct(*operation, randomNode[n].PublicKey)
		if err != nil {
			Log.Fatal("GetFile - Encrypting Data Failed")
		}

		opBox := OpBox{aesKey, encryptedOperation}
		operation.Next = opBox

		//operation.Next = encryptedOperation
		//operation.Op = ROUTE
	}
}
