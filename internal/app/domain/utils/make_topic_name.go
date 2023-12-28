package utils

import (
	"fmt"
	"watcher/internal/app/domain/entity"
)

func MakeBlocksTopic(appName string, blockchain entity.Blockchain, status entity.BlockStatus) string {
	return fmt.Sprint(appName, ".", entity.BlocksTopic, "-", blockchain, "-", status)
}
