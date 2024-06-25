package main

import (
	"github.com/go-grain/gin-layout/cmd/gen/sys"
)

// 第一次拉取项目,请先运行该文件,否则 repo层有依赖找不到
func main() {
	sys.Gen()
}
