package seafile

import (
	"path"
	"strings"
)

func buildObjectPaths(parentPath, parentPathInLibrary, entryPath, entryName string, recursive bool) (string, string) {
	var filePath, filePathInLibrary string
	if recursive {
		// In recursive mode, paths are built from DirEntry (+ a starting point)
		entryPath := strings.TrimPrefix(entryPath, "/")
		// If we're listing from some path inside the library (not the root)
		// there's already a path in parameter, which will also be included in the entry path
		entryPath = strings.TrimPrefix(entryPath, parentPathInLibrary)
		entryPath = strings.TrimPrefix(entryPath, "/")

		filePath = path.Join(parentPath, entryPath, entryName)
		filePathInLibrary = path.Join(parentPathInLibrary, entryPath, entryName)
	} else {
		// In non-recursive mode, paths are build from the parameters
		filePath = path.Join(parentPath, entryName)
		filePathInLibrary = path.Join(parentPathInLibrary, entryName)
	}
	return filePath, filePathInLibrary
}
