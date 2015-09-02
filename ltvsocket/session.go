package ltvsocket

import (
	"github.com/davyxu/cellnet"
	"log"
)

func SpawnSession(stream cellnet.PacketStream, createType SessionCreateType, callback func(interface{})) cellnet.CellID {

	recvCell := cellnet.Spawn(callback)

	// 发送线程
	sendCell := cellnet.Spawn(func(sendev interface{}) {

		if pkt, ok := sendev.(*cellnet.Packet); ok {
			stream.Write(pkt)
		} else {

			if config.SocketLog {
				log.Println("[ltvsocket] write require *cellnet.Packet type")
			}
		}

	})

	// 接收线程
	go func() {
		var err error
		var pkt *cellnet.Packet

		cellnet.LocalPost(recvCell, SocketNewSession{Session: sendCell, Type: createType})

		for {

			// 从Socket读取封包并转为ltv格式
			pkt, err = stream.Read()

			if err != nil {

				cellnet.Send(recvCell, SocketClose{Session: sendCell, Err: err})
				break
			}

			cellnet.LocalPost(recvCell, SocketData{Session: sendCell, Packet: pkt})

		}

	}()

	return recvCell
}
