// We can manually send "finish-it" message using "quit" channel.

package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/pborman/uuid"
)

// we need this only for visibility
// so we can have different colors for different sensors in command line
var (
	sensorConfig = map[string]string{
		uuid.New(): "\033[31m",
		uuid.New(): "\033[32m",
		uuid.New(): "\033[33m",
		uuid.New(): "\033[34m",
	}
)

type TempStat struct {
	SensorID string
	Temp     float32
}

func (t TempStat) String() string {
	return fmt.Sprintf("%sSensor ID: %s,\tSensor temperature value: %f", sensorConfig[t.SensorID], t.SensorID, t.Temp)
}

func fanIn(inputs [4]<-chan TempStat) <-chan TempStat {
	c := make(chan TempStat)
	go func() {
		for {
			select {
			case s := <-inputs[0]:
				c <- s
			case s := <-inputs[1]:
				c <- s
			case s := <-inputs[2]:
				c <- s
			case s := <-inputs[3]:
				c <- s
			}
		}
	}()

	return c
}

func main() {
	var sensors [4]<-chan TempStat

	quit := make(chan bool)

	i := 0
	for sensorID := range sensorConfig {
		sensors[i] = StatListen(sensorID, quit)
		i++
	}

	statsFanIn := fanIn(sensors)
	for i := 0; i < 5; i++ {
		s := <-statsFanIn
		fmt.Println(s)
	}

	quit <- true

	fmt.Println("Finished gathering sensor data.")
}

func StatListen(sensorID string, quit chan bool) <-chan TempStat { // Returns receive-only chan of TempStat.
	c := make(chan TempStat)
	go func() {
		for {
			select {
			case c <- getSensorTempStat(sensorID):
				// do something if you need to
			case <-quit:
				// some cleanup
				return
			}
		}
	}()
	return c
}

func getSensorTempStat(id string) TempStat {
	// simulate network latency or readiness of sensor
	time.Sleep(time.Duration(rand.Intn(2500)) * time.Millisecond)

	return TempStat{
		SensorID: id,
		// simulate temp difference
		Temp: float32((rand.Int63n(10) + 30)) * (rand.Float32() + 0.5),
	}
}
