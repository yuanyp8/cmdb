package version_test

import (
	"fmt"
	"github.com/yuanyp8/cmdb/version"
	"testing"
)

func TestFullVersion(t *testing.T) {
	fmt.Println(version.FullVersion())
}

func TestShortVersion(t *testing.T) {
	fmt.Println(version.ShortVersion())
}
