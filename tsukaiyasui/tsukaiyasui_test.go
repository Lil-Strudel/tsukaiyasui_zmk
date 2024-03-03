package tsukaiyasui

import (
	"fmt"
	"testing"
)

func Test(t *testing.T) {
	res := GenerateKeymap("corne")
	fmt.Println(res)

	GenerateBuildMatrix("corne", "nice_nano_v2")
}
