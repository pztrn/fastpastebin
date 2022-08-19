package helpers

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var (
	// ErrHelperError indicates that error came from helper.
	ErrHelperError = errors.New("helper")

	// ErrHelperPathNormalizationError indicates that error came from path normalization helper.
	ErrHelperPathNormalizationError = errors.New("path normalization")
)

// NormalizePath normalizes passed path:
// * Path will be absolute.
// * Symlinks will be resolved.
// * Support for tilde (~) for home path.
func NormalizePath(path string) (string, error) {
	// Replace possible tilde in the beginning (and only beginning!) of data path.
	if strings.HasPrefix(path, "~") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("%s: %s: %w", ErrHelperError, ErrHelperPathNormalizationError, err)
		}

		path = strings.Replace(path, "~", homeDir, 1)
	}

	// Normalize path.
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("%s: %s: %w", ErrHelperError, ErrHelperPathNormalizationError, err)
	}

	return absPath, nil
}
