package main

import (
	"fmt"
	"math/rand"
)

func main()  {
	//channelPipeline()
	//takeRepeat()
	repeatFn()
}

func channelPipeline()  {
	generator := func(done <-chan interface{}, intergers ...int) <-chan int {
		intStream := make(chan int, len(intergers))
		go func() {
			defer close(intStream)
			for _, i := range intergers {
				select {
				case <-done:
					return
				case intStream <- i:
				}
			}
		}()
		return intStream
	}

	multiply := func(done <-chan interface{}, intStream <-chan int, multiplier int) <-chan int {
		multipledStream := make(chan int)
		go func() {
			defer close(multipledStream)
			for i := range intStream {
				select {
				case <-done:
					return
				case multipledStream <- i*multiplier:
				}
			}
		}()
		return multipledStream
	}

	add := func(done <-chan interface{}, intStream <-chan int, additive int) <-chan int {
		addedStream := make(chan int)
		go func() {
			defer close(addedStream)
			for i := range intStream {
				select {
				case <-done:
					return
				case addedStream <- i+additive:
				}
			}
		}()
		return addedStream
	}

	done := make(chan interface{})
	defer close(done)

	intStream := generator(done, 1,2,3,4)
	fmt.Println(intStream)
	pipeline := multiply(done, add(done, multiply(done, intStream,2), 1), 2)

	for v := range pipeline {
		fmt.Println(v)
	}
}

func takeRepeat()  {
	done := make(chan interface{})
	defer close(done)

	// 渡した値をやめと送るまで無限に送り続ける関数
	repeat := func(done <-chan interface{}, values ...interface{}) <-chan interface{} {
		valueStream := make(chan interface{})
		go func() {
			defer close(valueStream)
			for  {
				for _, v := range values {
					select {
					case <-done:
						return
					case valueStream <- v:
					}
				}
			}
		}()
		return valueStream
	}

	// 受け取ったvalueStreamから num 個の要素だけ取り出して終了する。
	take := func(done <-chan interface{}, valueStream <-chan interface{}, num int) <-chan interface{} {
		takeStream := make(chan interface{})
		go func() {
			defer close(takeStream)
			for i := 0; i < num; i++ {
				select {
				case <-done:
					return
				case takeStream <- <- valueStream:

				}
			}
		}()
		return takeStream
	}

	for num := range take(done, repeat(done, 1), 10){
		fmt.Printf("%v ", num)
	}
}

func repeatFn()  {
	repeatFn := func(done <-chan interface{}, fn func() interface{}) <-chan interface{} {
		valueStream := make(chan interface{})
		go func() {
			defer close(valueStream)
			for  {
				select {
				case <-done:
					return
				case valueStream <- fn():
				}
			}
		}()
		return valueStream
	}

	take := func(done <-chan interface{}, valueStream <-chan interface{}, num int) <-chan interface{} {
		takeStream := make(chan interface{})
		go func() {
			defer close(takeStream)
			for i := 0; i < num; i++ {
				select {
				case <-done:
					return
				case takeStream <- <- valueStream:

				}
			}
		}()
		return takeStream
	}

	done := make(chan interface{})
	defer close(done)

	rand := func() interface{} { return rand.Int() }

	for num := range take(done, repeatFn(done, rand), 10) {
		fmt.Println(num)
	}
}
