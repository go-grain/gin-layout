package sys

import (
	sysModel "github.com/go-grain/gin-layout/model/system"
	"gorm.io/gen"
)

func Gen() {
	xGen(
		"internal/repo/system/query",
		sysModel.SysRole{},
		sysModel.SysUser{},
	)
}

func xGen(outPath string, model ...interface{}) {
	g := gen.NewGenerator(gen.Config{
		OutPath: outPath,
		Mode:    gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface,
	})

	g.ApplyBasic(
		model...,
	)
	g.Execute()
}
