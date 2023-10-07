package test

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
)

func Test_chan(t *testing.T) {
	numInputData := 10
	inputChan := make(chan int, numInputData)
	// 向输入channel中写入10000个数据
	for i := 0; i < numInputData; i++ {
		index := i
		inputChan <- index
	}
	close(inputChan)
	for index := range inputChan {
		fmt.Println(index)
	}
}

func Test_uuid(t *testing.T) {
	u := uuid.New()
	fmt.Println(u.String())
}
