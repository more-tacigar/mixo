// ==================================================
// Copyright (c) 2016 tacigar. All rights reserved.
// ==================================================

package mixo

func calculateRelativePath(path1, path2 string) string {
	rpath1 := []rune(path1)
	rpath2 := []rune(path2)
	var i = 0
	for ; i < len(rpath1); i++ {
		if rpath1[i] != rpath2[i] {
			return ""
		}
	}
	return string(path2[i:])
}
