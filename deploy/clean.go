package deploy

import (
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/errgo.v1"
)

func cleanDestination(dir string) error {
	needRemove, err := canClean(dir)
	if err != nil {
		return errgo.Mask(err)
	}
	for _, p := range needRemove {
		if err := os.RemoveAll(p); err != nil {
			return errgo.Mask(err)
		}
	}
	return nil
}

var allowed = map[string]bool{
	"bin":              true,
	"assets":           true,
	"compile":          true,
	"config.yaml":      true,
	"dependencies.tsv": true,
	"hooks":            true,
	"metadata.yaml":    true,
	"pkg":              true, // This allows us to test the compile scripts in the charm dir.
	"README.md":        true,
	"revision":         true,
	"src":              true,
}

func canClean(dir string) (needRemove []string, err error) {
	infos, err := ioutil.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, errgo.Mask(err)
	}
	var toRemove []string
	for _, info := range infos {
		if info.Name()[0] == '.' {
			continue
		}
		if !allowed[info.Name()] {
			return nil, errgo.Newf("unexpected file %q found in %s", info.Name(), dir)
		}
		path := filepath.Join(dir, info.Name())
		if strings.HasSuffix(path, ".yaml") && !autogenerated(path) {
			return nil, errgo.Newf("non-autogenerated file %q", path)
		}
		toRemove = append(toRemove, path)
	}
	return toRemove, nil
}

func autogenerated(path string) bool {
	f, err := os.Open(path)
	if err != nil {
		return false
	}
	defer f.Close()
	buf := make([]byte, len(yamlAutogenComment))
	if _, err := io.ReadFull(f, buf); err != nil {
		return false
	}
	return bytes.Equal(buf, []byte(yamlAutogenComment))
}
