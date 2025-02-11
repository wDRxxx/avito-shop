package utils

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func ParseCustomDuration(input string) (time.Duration, error) {
	if strings.HasSuffix(input, "d") {
		daysStr := strings.TrimSuffix(input, "d")
		days, err := strconv.Atoi(daysStr)
		if err != nil {
			return 0, fmt.Errorf("не удалось распарсить дни: %v", err)
		}

		return time.Duration(days) * 24 * time.Hour, nil
	}

	return time.ParseDuration(input)
}
