// ==================================================
// Copyright (c) 2016 tacigar. All rights reserved.
// ==================================================

package mixo

import (
	"testing"
)

func TestCalculateRelativePath(t *testing.T) {
	path1 := "/auth/"
	path2 := "/auth/test"
	if s := calculateRelativePath(path1, path2); s != "test" {
		t.Errorf("error: path2 - path1 must be \"test\": %s", s)
	}
}
