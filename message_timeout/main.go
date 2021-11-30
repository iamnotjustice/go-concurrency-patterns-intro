// If we need to make sure our messages arrive in expected timeframe we can use timeouts!
// It's basically the same example as seelct_fan_in, but with message timeout.

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
			// if no other channel is ready in 1 second - we shut everything down
			case <-time.After(1 * time.Second):
				fmt.Println("Waited too long!")
				close(c)
				return
			}
		}
	}()

	return c
}

func main() {
	var sensors [4]<-chan TempStat

	i := 0
	for sensorID := range sensorConfig {
		sensors[i] = StatListen(sensorID)
		i++
	}

	statsFanIn := fanIn(sensors)

	for v := range statsFanIn {
		fmt.Printf("%v\n", v)
	}

	fmt.Println("Finished gathering sensor data.")
}

func StatListen(sensorID string) <-chan TempStat { // Returns receive-only chan of TempStat.
	c := make(chan TempStat)

	go func() {
		for {
			c <- getSensorTempStat(sensorID)
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
