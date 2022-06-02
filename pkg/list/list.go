package list

import (
	"fmt"
	"os"
	"path/filepath"
)

func Run(home string) ([]string, error) {
	var list []string
	dirs, err := os.ReadDir(home)
	if err != nil {
		return nil, fmt.Errorf("failed to read %s: %v", home, err)
	}
	for _, dir := range dirs {
		if !dir.IsDir() {
			continue
		}
		groupName := dir.Name()
		groupPath := filepath.Join(home, dir.Name())
		projs, err := os.ReadDir(groupPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read %s: %v", groupPath, err)
		}
		for _, proj := range projs {
			item := fmt.Sprintf("%s/%s", groupName, proj.Name())
			list = append(list, item)
		}
	}
	return list, nil
}
