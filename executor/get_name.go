package executor

import (
	"fmt"
	"time"
)

func getName() string {
	return fmt.Sprintf("filewalker-%d", time.Now().UnixNano())
}
