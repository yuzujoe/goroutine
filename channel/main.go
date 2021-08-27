package main

import (
	"bytes"
	"fmt"
	"os"
	"sync"
)

func main()  {
	//multipleGoroutineClose()
	//bufferChannel()
	chanOwner()
}

func deadLock()  {
	stringStream := make(chan string)
	go func() {
		if 0 != 1 {
			return
		}
		stringStream <- "Hello channels!"
	}()
	fmt.Println(<-stringStream)
}

func options()  {
	stringStream := make(chan string)
	go func() {
		stringStream <- "Hello channels"
	}()
	salutation, ok :=  <-stringStream
	fmt.Printf("(%v), %v", ok, salutation)
}

func chanClose()  {
	intStream := make(chan int)
	go func() {
		defer close(intStream)
		for i := 1; i <= 5; i++ {
			intStream <- i
		}
	}()

	for integer := range intStream {
		fmt.Printf("%v ", integer)
	}
}

func multipleGoroutineClose()  {
	begin := make(chan interface{})
	var wg sync.WaitGroup
	for i := 0; i <= 4; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			<-begin
			fmt.Printf("%v has begun\n", i)
		}(i)
	}

	fmt.Println("Unblocking goroutines...")
	close(begin)
	wg.Wait()
}

func bufferChannel()  {
	var stdoutBuff bytes.Buffer
	defer stdoutBuff.WriteTo(os.Stdout)

	intStream := make(chan int, 4)
	go func() {
		defer close(intStream)
		defer fmt.Fprintln(&stdoutBuff, "Producer Done.")
		for i := 0; i < 5; i++ {
			fmt.Fprintf(&stdoutBuff, "Sending: %d\n", i)
			intStream <- i
		}
	}()

	for integer := range intStream{
		fmt.Fprintf(&stdoutBuff, "Received %v. \n", integer)
	}
}

func chanOwner()  {
	chanOwner := func() <-chan int {
		resultStream := make(chan int, 5)
		go func() {
			defer close(resultStream)
			for i := 0; i <= 5; i++ {
				resultStream<-i
			}
		}()
		return resultStream
	}

	resultStream := chanOwner()
	for result := range resultStream {
		fmt.Printf("Receiveed: %d\n", result)
	}
	fmt.Println("Done receiving!")
}
