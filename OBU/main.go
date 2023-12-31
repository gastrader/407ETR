// sits in cars and sends out GPS coordinates at intervals
package main

import (
	"log"

	"math/rand"
	"time"

	"github.com/gastrader/407ETR/types"
	"github.com/gorilla/websocket"
)

const wsEndpoint = "ws://127.0.0.1:30000/ws"

var sendInterval = time.Second * 5


func genCoord() float64 {
	n := float64(rand.Intn(100) + 1)
	f := rand.Float64()
	return n + f
}

func genLatLong() (float64, float64) {
	return genCoord(), genCoord()
}


func gnereateOBUIDS(n int) []int {
	ids := make([]int, n)
	for i := 0; i < n; i++ {
		ids[i]=rand.Intn(999999)
	}
	return ids
}

func main() {
	obuIDS := gnereateOBUIDS(3)
	conn, _, err := websocket.DefaultDialer.Dial(wsEndpoint, nil)
	if err != nil{
		log.Fatal(err)
	}
	for {
		for i := 0; i<len(obuIDS); i++ {
			lat, long := genLatLong()
			data := types.OBUData{
				OBUID: obuIDS[i],
				Lat: lat,
				Long: long,
			}
			if err := conn.WriteJSON(data); err != nil{
				log.Fatal(err)
			}
		}
		time.Sleep(sendInterval)
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
