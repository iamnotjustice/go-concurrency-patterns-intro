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
		uuid.New(): "\033[35m",
		uuid.New(): "\033[36m",
		uuid.New(): "\033[37m",
	}
)

type TempStat struct {
	SensorID string
	Temp     float32
}

func (t TempStat) String() string {
	return fmt.Sprintf("%sSensor ID: %s,\tSensor temperature value: %f", sensorConfig[t.SensorID], t.SensorID, t.Temp)
}

func fanIn(inputs ...<-chan TempStat) <-chan TempStat {
	c := make(chan TempStat)
	for _, v := range inputs {
		v := v
		go func() {
			for {
				c <- <-v
			}
		}()
	}

	return c
}

func main() {
	// we got 5 "sensor" channels
	var sensors []<-chan TempStat

	for sensorID := range sensorConfig {
		// we start listening to sensor data here and append it to the list
		sensors = append(sensors, StatListen(sensorID))
	}

	// we gather all data from all sensors using one fanIn channel
	statsFanIn := fanIn(sensors...)

	for v := range statsFanIn {
		// here we simulate collecting data
		// or sending it over the network
		// or adding values to some algorithm that calculates stuff
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
