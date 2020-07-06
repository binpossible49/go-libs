package main

import (
	"fmt"

	log "github.com/binpossible49/go-libs/log"
	"github.com/binpossible49/go-libs/registry"

	"go.uber.org/zap"
)

func main() {
	fmt.Println("Hello")
	log.InitZap("go-libs", "D", map[string]string{})
	registry.BuildDIContainer()
	err := log.InitZap("Lib", "D", map[string]string{
		"cardno": "(?P<FIRST>[0-9]{6})(?P<MASK>[0-9]*)(?P<LAST>[0-9]{4})",
	})
	fmt.Println(err)

	zap.S().Info(Test{Name: "Test", CardNo: "12345678912345678"})
}

type Test struct {
	Name   string
	CardNo string
}
