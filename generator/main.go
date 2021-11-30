package main

import (
	"fmt"
	"math/rand"
	"time"
)

func boring(msg string) <-chan string { // Возвращает receive-only канал строк.
	c := make(chan string)
	go func() { // Запускаем горутину внутри функции.
		for i := 0; ; i++ {
			c <- fmt.Sprintf("%s %d", msg, i)
			time.Sleep(time.Duration(rand.Intn(1e3)) * time.Millisecond)
		}
	}()
	return c // Возвращаем канал.
}

func main() {
	joe := boring("Tom Henderson")
	ann := boring("Ann Perkins")
	for i := 0; i < 5; i++ {
		fmt.Println(<-joe)
		fmt.Println(<-ann)
	}
	fmt.Println("You're boring; I'm leaving.")
}
