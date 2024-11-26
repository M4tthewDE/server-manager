package memory

import (
	"os"
	"strings"
)

type Memory struct {
	Total     string
	Free      string
	Available string
}

func FetchMemory() (*Memory, error) {
	data, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return nil, err
	}

	total, err := parseTotal(string(data))
	if err != nil {
		return nil, err
	}

	free, err := parseFree(string(data))
	if err != nil {
		return nil, err
	}

	available, err := parseAvailable(string(data))
	if err != nil {
		return nil, err
	}

	return &Memory{
		Total:     total,
		Free:      free,
		Available: available,
	}, nil
}

func parseTotal(text string) (string, error) {
	line := strings.Split(text, "\n")[0]
	total := strings.Split(line, ":")[1]
	total = strings.TrimSpace(total)
	return total, nil
}

func parseFree(text string) (string, error) {
	line := strings.Split(text, "\n")[1]
	total := strings.Split(line, ":")[1]
	total = strings.TrimSpace(total)
	return total, nil
}

func parseAvailable(text string) (string, error) {
	line := strings.Split(text, "\n")[2]
	total := strings.Split(line, ":")[1]
	total = strings.TrimSpace(total)
	return total, nil
}
