package async

import (
	"go.uber.org/zap"
	"watcher/configs"
	"watcher/internal/app/domain/entity"
	"watcher/internal/app/usecase"
)

type WatchProcess struct {
	logger              *zap.Logger
	cfg                 configs.Config
	confirmationUsecase *usecase.BlockConfirmationUsecase
	tickUsecase         *usecase.BlockHandlingUsecase
	confirmationsChan   chan entity.BlockHeight
}

func NewWatcher(
	logger *zap.Logger,
	cfg configs.Config,
	confirmationUsecase *usecase.BlockConfirmationUsecase,
	tickUsecase *usecase.BlockHandlingUsecase,
) *WatchProcess {
	return &WatchProcess{
		logger:              logger,
		cfg:                 cfg,
		confirmationUsecase: confirmationUsecase,
		tickUsecase:         tickUsecase,
		confirmationsChan:   make(chan entity.BlockHeight),
	}
}
