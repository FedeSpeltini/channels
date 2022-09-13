package main

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
)

func main() {
	wg := &sync.WaitGroup{}
	IDsChan := make(chan string)
	FakeIDsChan := make(chan string)
	closedChan := make(chan int)

	wg.Add(3)
	go generateIds(wg, IDsChan, closedChan)
	go generateFakeIds(wg, FakeIDsChan, closedChan)
	go logIDs(wg, IDsChan, FakeIDsChan, closedChan)

	wg.Wait()
}

func generateFakeIds(wg *sync.WaitGroup, fakeIDsChan chan<- string, closedChan chan<- int) {
	for i := 0; i < 100; i++ {
		id := uuid.New()
		fakeIDsChan <- fmt.Sprintf("%d. %s", i+1, id.String())
	}
	close(fakeIDsChan)
	closedChan <- 1
	wg.Done()

}

func generateIds(wg *sync.WaitGroup, idsChan chan<- string, closedChan chan<- int) {
	for i := 0; i < 100; i++ {
		id := uuid.New()
		idsChan <- fmt.Sprintf("%d. %s", i+1, id.String())
	}
	close(idsChan)
	closedChan <- 1
	wg.Done()

}

func logIDs(wg *sync.WaitGroup, idsChan <-chan string, fakeIDsChan <-chan string, closedChan chan int) {
	closedCounter := 0

	for {
		select {
		case id, ok := <-idsChan:
			if ok {
				fmt.Println("ID:" + id)
			}
		case id, ok := <-fakeIDsChan:
			if ok {
				fmt.Println("FAKE ID:" + id)
			}
		case count, ok := <-closedChan:
			if ok {
				closedCounter += count
			}
		}
		if closedCounter == 2 {
			close(closedChan)
			break
		}
	}

	wg.Done()
}
