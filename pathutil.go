package main

import (
	"os"
	"path/filepath"
)

// resolvedMkdirAll creates a directory and all necessary parents,
// resolving firmlinks and symlinks on the existing portion of the
// path before calling os.MkdirAll. This avoids ENOTDIR failures
// on macOS Big Sur+ where /Volumes is a firmlink to a synthetic
// data volume and os.MkdirAll's stat-based traversal gets confused.
func resolvedMkdirAll(path string, perm os.FileMode) error {
	// Walk up from path to find the deepest existing ancestor.
	existing := path
	var tail []string
	for {
		_, err := os.Lstat(existing)
		if err == nil {
			break
		}
		tail = append([]string{filepath.Base(existing)}, tail...)
		parent := filepath.Dir(existing)
		if parent == existing {
			break // reached root
		}
		existing = parent
	}

	// Resolve firmlinks/symlinks on the existing portion.
	resolved, err := filepath.EvalSymlinks(existing)
	if err != nil {
		resolved = existing // fall back to original
	}

	// Reconstruct the full path with the resolved base.
	full := resolved
	for _, component := range tail {
		full = filepath.Join(full, component)
	}

	return os.MkdirAll(full, perm)
}
