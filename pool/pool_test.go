package pool

import (
	"log"
	"math/rand"
	"os"
	"testing"
	"time"
)

var (
	poolSize  = 100
	taskCount = 100000
	pip       = make(chan int, 100)
)

func produce() {
	pip <- rand.Intn(100)
}

func consume() {
	log.Printf("consume %d\n", <-pip)
}

func TestProduce(t *testing.T) {
	f, err := os.OpenFile("./produce.log", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(f)

	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(f)

	producerPool := New(poolSize)
	go func() {
		for i := 0; i < taskCount; i++ {
			producerPool.Put(produce)
		}
		log.Println("produce ended")
	}()

	go func() {
		for i := 0; i < taskCount; i++ {
			consume()
		}
		log.Println("consume ended")
	}()

	time.Sleep(time.Second * 5)
	producerPool.Stop()
}

func TestConsume(t *testing.T) {
	f, err := os.OpenFile("./consume.log", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(f)

	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(f)

	go func() {
		for i := 0; i < taskCount; i++ {
			pip <- i
		}
		log.Println("produce ended")
	}()

	consumerPool := New(poolSize)
	go func() {
		for i := 0; i < taskCount; i++ {
			consumerPool.Put(consume)
		}
		log.Println("consume ended")
	}()

	time.Sleep(time.Second * 5)
	consumerPool.Stop()
}
