package logic

import (
	"sync"
)

func process(numWorkers, numInputData int, handler func(index int)) {
	inputChan := make(chan int, numInputData)
	// 向输入channel中写入10000个数据
	for i := 0; i < numInputData; i++ {
		index := i
		inputChan <- index
	}
	close(inputChan)
	var wg sync.WaitGroup
	var worker = func(in <-chan int) {
		defer wg.Done()
		for index := range in {
			handler(index)
		}
	}
	wg.Add(numWorkers)
	// 启动10个worker goroutine并从输入channel中读取数据进行处理
	for i := 0; i < numWorkers; i++ {
		go worker(inputChan)
	}
	wg.Wait()
}
