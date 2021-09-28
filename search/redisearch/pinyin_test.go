package main

import (
	"fmt"
	"testing"
)

func TestPinYin(t *testing.T) {
	name := "基础设施/环境空调/温湿度检测/温湿度/温湿度/TH03/TH03"
	fmt.Println(PinYin(name))
}
