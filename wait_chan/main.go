// Simplified sensor example to show how we can send a channel as part of the channel data,
// effectively giving the receiver a handler using which it can control the flow of the goroutine.

package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Message struct {
	stat TempStat
	wait chan int
}

type TempStat struct {
	SensorID string
	Temp     float32
}

func (t TempStat) String() string {
	return fmt.Sprintf("Sensor ID: %s,\tSensor temperature value: %f", t.SensorID, t.Temp)
}

func fanIn(input1, input2 <-chan Message) <-chan Message {
	c := make(chan Message)

	go func() {
		for {
			c <- <-input1
		}
	}()

	go func() {
		for {
			c <- <-input2
		}
	}()

	return c
}

func main() {
	statsFanIn := fanIn(StatListen("1"), StatListen("2"))
	for {
		msg1 := <-statsFanIn
		fmt.Printf("%v\n", msg1.stat)
		msg2 := <-statsFanIn
		fmt.Printf("%v\n", msg2.stat)

		<-msg1.wait
		<-msg2.wait
	}
}

func StatListen(sensorID string) <-chan Message { // Returns receive-only chan of TempStat.
	c := make(chan Message)
	waitingStep := make(chan int)
	go func() {
		for {
			c <- Message{getSensorTempStat(sensorID), waitingStep}
			time.Sleep(time.Second)
			waitingStep <- 1 // we block this goroutine until the next channel operation
		}
	}()
	return c
}

func getSensorTempStat(id string) TempStat {
	return TempStat{
		SensorID: id,
		// simulate temp difference
		Temp: float32((rand.Int63n(10) + 30)) * (rand.Float32() + 0.5),
	}
}
