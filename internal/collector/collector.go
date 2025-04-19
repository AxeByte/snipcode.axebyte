package collector

import (
  "io/fs"
  "path/filepath"
  "sort"

  "github.com/bmatcuk/doublestar/v4"
)

// Collect returns files not matching ignore patterns
func Collect(root string, patterns []string) ([]string, error) {
  var files []string
  err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
    if err != nil || d.IsDir() || path == ".grepattern.yaml" {
      return err
    }
    for _, pat := range patterns {
      if ok, _ := doublestar.PathMatch(pat, path); ok {
        return nil
      }
    }
    files = append(files, path)
    return nil
  })
  sort.Strings(files)
  return files, err
}
