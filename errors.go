// ==================================================
// Copyright (c) 2016 tacigar. All rights reserved.
// ==================================================

package mixo

import (
	"errors"
)

var (
	ErrAddingRegisteredRoute = errors.New("error : cannot register the route registered already")
)
