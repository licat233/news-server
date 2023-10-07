package utils

import (
	"math/rand"
	"time"

	"github.com/google/uuid"
)

func RandNumbers(max int64, count int) []int {
	// 设置随机种子，确保每次运行结果不同
	randCli := rand.New(rand.NewSource(time.Now().UnixNano()))

	nums := make([]int, 1500) // 创建包含1500个数的切片
	for i := 0; i < 1500; i++ {
		nums[i] = i + 1 // 将数值填入切片中
	}

	randCli.Shuffle(len(nums), func(i, j int) {
		nums[i], nums[j] = nums[j], nums[i] // 随机打乱切片顺序
	})

	// 输出前10个随机数
	return nums[:count]
}

func UUID() string {
	return uuid.New().String()
}
