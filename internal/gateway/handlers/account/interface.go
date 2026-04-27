package account

import (
	"github.com/vyolayer/vyolayer/internal/gateway/service"
	"github.com/vyolayer/vyolayer/pkg/jwt"
	"github.com/vyolayer/vyolayer/pkg/logger"
	accountV1 "github.com/vyolayer/vyolayer/proto/account/v1"
)

// AccountHandler manages HTTP requests related to accounts
type AccountHandler struct {
	client     accountV1.AccountServiceClient
	cookieSv   *service.AccountTokenService
	accountJWT jwt.AccountJWT
	logger     *logger.AppLogger
}
