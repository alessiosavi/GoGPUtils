package mathutils

import (
	"testing"
)

func TestConvertSize(t *testing.T) {
	t.Log(ConvertSize(1024, "KB"))
	t.Log(ConvertSize(1000000, "MB"))
	t.Log(ConvertSize(1024, "GB"))
}
