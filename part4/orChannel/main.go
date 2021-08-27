package main

import (
	"fmt"
	"time"
)

func simpleOrChannel()  {
	var or func(channels ...<-chan interface{}) <-chan interface{}
	// 1 関数 or を定義して1つのチャネルを返す
	or = func(channels ...<-chan interface{}) <-chan interface{} {
		switch len(channels) {
		// 2. 停止条件を定める、合成チャネルを作ることは想定していない
		case 0:
			return nil
		// 3. スライスが1つしか要素を持っていなかった場合はその要素を返す
		case 1:
			return channels[0]
		}

		orDone := make(chan interface{})
		// 4. ゴルーチンを作成してブロックすることなく作ったチャネルにメッセージを受け取れる様にしている。
		go func() {
			defer close(orDone)

			switch len(channels) {
			// 5. or への再帰の呼び出しは少なくとも2つのチャネルを持っているのでその為の条件を作成する
			case 2:
				select {
				case <-channels[0]:
				case <-channels[1]:
				}
			default:
				select {
				case <-channels[0]:
				case <-channels[1]:
				case <-channels[2]:
				case <-or(append(channels[3:], orDone)...):
				}
			}
		}()
		return orDone
	}
	sig := func(after time.Duration) <-chan interface{} {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	start := time.Now()
	<-or(
		sig(2*time.Hour),
		sig(5*time.Minute),
		sig(1*time.Second),
		sig(1*time.Hour),
		sig(1*time.Minute),
	)
	fmt.Printf("done after %v", time.Since(start))
}


func main()  {
 	simpleOrChannel()
}
