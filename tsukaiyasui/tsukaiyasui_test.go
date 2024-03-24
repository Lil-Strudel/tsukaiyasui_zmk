package tsukaiyasui

import (
	"testing"
)

func Test(t *testing.T) {
	GenerateKeymap("corne")
	GenerateBuildMatrix("corne", "nice_nano_v2")
}
