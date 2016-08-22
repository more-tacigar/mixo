// ==================================================
// Copyright (c) 2016 tacigar. All rights reserved.
// ==================================================

package mixo

import (
	"testing"
)

func TestDivideSessionFileLine(t *testing.T) {
	test := "key=value"
	if s1, s2 := divideSessionFileLine(test); s1 != "key" || s2 != "value" {
		t.Errorf("must be \"key\" and \"value\": %s, %s", s1, s2)
	}
}

func TestSessionStart(t *testing.T) {
	context := &Context{}
	h := SessionStart("session", 10)
	h(context)

	if _, b := context.GetMetadata(sessionDefaultKey); !b {
		t.Errorf("context must have session manager")
	}
}

func TestGetSessionManager(t *testing.T) {
	context := &Context{}
	h := SessionStart("session", 10)
	h(context)

	if GetSessionManager(context) == nil {
		t.Errorf("context must have session manager")
	}
}
