// Same thing, just like in fan_in_sensors we gather data as it comes,
// but here we use select to control the flow of the fanIn function.

// Notice how we need to know exactly how many inputs we have to create a proper select statement!

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

// cases need to be hardcoded, hence we don't use a variadic parameters here
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
	time.Sleep(time.Duration(rand.Intn(10)) * time.Second)

	return TempStat{
		SensorID: id,
		// simulate temp difference
		Temp: float32((rand.Int63n(10) + 30)) * (rand.Float32() + 0.5),
	}
}
