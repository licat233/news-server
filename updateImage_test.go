package main

import (
	"context"
	"log"
	_ "news/global"
	"news/logic"
	"testing"
)

func Test_UpdateImages(t *testing.B) {
	err := logic.UpdateImageHandler(context.Background())
	log.Fatalln(err)
}
