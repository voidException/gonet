package protos

import (
	"hub/ranklist"
	"misc/packet"
)

import (
	"fmt"
	"log"
	"runtime"
)

func HandleRequest(hostid int32, p []byte) []byte {
	defer _HandleError()

	reader := packet.Reader(p)
	b, err := reader.ReadU16()

	if err != nil {
		log.Println("read protocol error")
	}

	handle := ProtoHandler[b]
	if handle != nil {
		ret, err := handle(hostid, reader)
		fmt.Println(ret)
		if err == nil {
			return ret
		}
	}

	return nil
}

func _login_req(hostid int32, pkt *packet.Packet) ([]byte, error) {
	tbl, _ := pktread_id(pkt)
	ret := intresult{F_v: 0}

	if ranklist.Login(tbl.F_id, hostid) {
		ret.F_v = 1
	}

	return packet.Pack(Code["login_ack"], ret, nil), nil
}

func _logout_req(hostid int32, pkt *packet.Packet) ([]byte, error) {
	tbl, _ := pktread_id(pkt)
	ret := intresult{F_v: 0}

	if ranklist.Login(tbl.F_id, hostid) {
		ret.F_v = 1
	}

	return packet.Pack(Code["logout_ack"], ret, nil), nil
}

func _changescore_req(hostid int32, pkt *packet.Packet) ([]byte, error) {
	tbl, _ := pktread_changescore(pkt)
	ret := intresult{F_v: 0}

	if ranklist.ChangeScore(tbl.F_id, tbl.F_oldscore, tbl.F_newscore) {
		ret.F_v = 1
	}

	return packet.Pack(Code["changescore_ack"], ret, nil), nil
}

func _getlist_req(hostid int32, pkt *packet.Packet) ([]byte, error) {
	tbl, _ := pktread_getlist(pkt)
	ret := getlist_result{}

	ids, scores := ranklist.GetList(int(tbl.F_A), int(tbl.F_B))
	ret.F_items=make([]id_score,len(ids))

	for k := range ids {
		ret.F_items[k].F_id = ids[k]
		ret.F_items[k].F_score = scores[k]
	}

	return packet.Pack(Code["getlist_ack"], ret, nil), nil
}

func _raid_req(hostid int32, pkt *packet.Packet) ([]byte, error) {
	return nil, nil
}

func _protect_req(hostid int32, pkt *packet.Packet) ([]byte, error) {
	return nil, nil
}

func _unprotect_req(hostid int32, pkt *packet.Packet) ([]byte, error) {
	return nil, nil
}

func _free_req(hostid int32, pkt *packet.Packet) ([]byte, error) {
	return nil, nil
}

func _getstate_req(hostid int32, pkt *packet.Packet) ([]byte, error) {
	return nil, nil
}

func _getprotecttime_req(hostid int32, pkt *packet.Packet) ([]byte, error) {
	return nil, nil
}

func _getname_req(hostid int32, pkt *packet.Packet) ([]byte, error) {
	return nil, nil
}

func checkErr(err error) {
	if err != nil {
		funcName, file, line, ok := runtime.Caller(1)
		if ok {
			log.Printf("ERR:%v,[func:%v,file:%v,line:%v]\n", err, runtime.FuncForPC(funcName).Name(), file, line)
		}

		panic("error occured in protocol module")
	}
}

func _HandleError() {
	if x := recover(); x != nil {
		log.Printf("run time panic when processing HUB request: %v", x)
		for i := 0; i < 10; i++ {
			funcName, file, line, ok := runtime.Caller(i)
			if ok {
				log.Printf("frame %v:[func:%v,file:%v,line:%v]\n", i, runtime.FuncForPC(funcName).Name(), file, line)
			}
		}
	}
}
