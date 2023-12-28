package async

import (
	"go.uber.org/zap"
	"watcher/configs"
	"watcher/internal/app/usecase"
)

type WatchProcess struct {
	logger               *zap.Logger
	cfg                  configs.Config
	confirmationUsecase  *usecase.BlockConfirmationUsecase
	blockHandlingUsecase *usecase.BlockHandlingUsecase
	cleanupUsecase       *usecase.CleanupUsecase
}

func NewWatcher(
	logger *zap.Logger,
	cfg configs.Config,
	confirmationUsecase *usecase.BlockConfirmationUsecase,
	blockHandlingUsecase *usecase.BlockHandlingUsecase,
	cleanupUsecase *usecase.CleanupUsecase,
) *WatchProcess {
	return &WatchProcess{
		logger:               logger,
		cfg:                  cfg,
		confirmationUsecase:  confirmationUsecase,
		blockHandlingUsecase: blockHandlingUsecase,
		cleanupUsecase:       cleanupUsecase,
	}
}
