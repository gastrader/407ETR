package main

import (
	"fmt"
	"log"
	"net/http"

	// "github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gastrader/407ETR/types"
	"github.com/gorilla/websocket"
)

func main() {

	recv, err := NewDataReceiver()
	if err != nil{
		log.Fatal(err)
	}
	http.HandleFunc("/ws", recv.handleWS)
	fmt.Println("data receiver working")
	http.ListenAndServe(":30000", nil)
}

type DataReceiver struct {
	msgch chan types.OBUData
	conn  *websocket.Conn
	prod  DataProducer
}
func NewDataReceiver() (*DataReceiver, error){
	var (
		p DataProducer
		err error
		kafkaTopic = "obuData"
	)
	p, err = NewKafkaProducer(kafkaTopic)
	if err != nil {
		return nil, err
	}
	p = NewLogMiddleware(p)
	if err != nil {
		return nil, err
	}
	return &DataReceiver{
		msgch: make(chan types.OBUData, 128),
		prod: p,
	}, nil
}

func (dr *DataReceiver) produceData(data types.OBUData) error {
	return dr.prod.ProduceData(data)
}

func (dr *DataReceiver) handleWS(w http.ResponseWriter, r *http.Request) {
	u := websocket.Upgrader{
		ReadBufferSize:  1028,
		WriteBufferSize: 1028,
	}
	conn, err := u.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}
	dr.conn = conn

	go dr.wsReceiveLoop()
}

func (dr *DataReceiver) wsReceiveLoop() {
	fmt.Println("New OBU/Client Connected!")
	for {
		var data types.OBUData
		if err := dr.conn.ReadJSON(&data); err != nil {
			log.Println("read err: ", err)
			continue
		}
		fmt.Printf("received: %+v\n ", data)	
		if err := dr.produceData(data); err != nil {
			fmt.Println("kafka error: ", err)
		}
		
	}
}
