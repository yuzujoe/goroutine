package main

import (
	"fmt"
	"sync"
)

func main()  {
	calcPool()
}

func pool()  {
	myPool := &sync.Pool{
		New: func() interface{} {
			fmt.Println("Createing new instance.")
			return struct {}{}
		},
	}

	myPool.Get()
	instance := myPool.Get()
	myPool.Put(instance)
	myPool.Get()
}

func calcPool()  {
	var numCalcsCreated int
	calcPool := &sync.Pool{
		New: func() interface{} {
			numCalcsCreated += 1
			mem := make([]byte, 1024)
			return &mem
		},
	}

	calcPool.Put(calcPool.New())
	calcPool.Put(calcPool.New())
	calcPool.Put(calcPool.New())
	calcPool.Put(calcPool.New())

	const numWorks = 1024*1024
	var wg sync.WaitGroup
	wg.Add(numWorks)
	for i := numWorks; i > 0; i-- {
		go func() {
			defer wg.Done()

			mem := calcPool.Get().(*[]byte)
			defer calcPool.Put(mem)
		}()
	}

	wg.Wait()
	fmt.Printf("%d calculators ware created.", numCalcsCreated)
}


