package memory

import (
	"os"
	"strconv"
	"strings"

	"github.com/dustin/go-humanize"
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
	return parse(line)
}

func parseFree(text string) (string, error) {
	line := strings.Split(text, "\n")[1]
	return parse(line)
}

func parseAvailable(text string) (string, error) {
	line := strings.Split(text, "\n")[2]
	return parse(line)
}

func parse(line string) (string, error) {
	rightSide := strings.Split(line, ":")[1]
	rightSide = strings.TrimSpace(rightSide)
	rightSide = strings.Split(rightSide, " ")[0]
	kilobytes, err := strconv.Atoi(rightSide)
	if err != nil {
		return "", err
	}

	return humanize.SIWithDigits(float64(kilobytes)*1000, 3, "B"), nil
}
