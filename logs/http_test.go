package logs_test

import (
	"net/http"
	"testing"

	. "github.com/bbengfort/catena/logs"
	"github.com/stretchr/testify/require"
)

func TestStatusLevel(t *testing.T) {
	tt := []struct {
		status int
		level  uint8
	}{
		{http.StatusContinue, LevelInfo},
		{http.StatusSwitchingProtocols, LevelInfo},
		{http.StatusEarlyHints, LevelInfo},
		{http.StatusOK, LevelStatus},
		{http.StatusCreated, LevelStatus},
		{http.StatusAccepted, LevelStatus},
		{http.StatusNonAuthoritativeInfo, LevelStatus},
		{http.StatusNoContent, LevelStatus},
		{http.StatusResetContent, LevelStatus},
		{http.StatusPartialContent, LevelStatus},
		{http.StatusMultipleChoices, LevelStatus},
		{http.StatusMovedPermanently, LevelStatus},
		{http.StatusFound, LevelStatus},
		{http.StatusSeeOther, LevelStatus},
		{http.StatusNotModified, LevelStatus},
		{http.StatusTemporaryRedirect, LevelStatus},
		{http.StatusPermanentRedirect, LevelStatus},
		{http.StatusBadRequest, LevelWarn},
		{http.StatusUnauthorized, LevelWarn},
		{http.StatusPaymentRequired, LevelWarn},
		{http.StatusForbidden, LevelWarn},
		{http.StatusNotFound, LevelWarn},
		{http.StatusMethodNotAllowed, LevelWarn},
		{http.StatusNotAcceptable, LevelWarn},
		{http.StatusProxyAuthRequired, LevelWarn},
		{http.StatusRequestTimeout, LevelWarn},
		{http.StatusConflict, LevelWarn},
		{http.StatusGone, LevelWarn},
		{http.StatusInternalServerError, LevelWarn},
		{http.StatusNotImplemented, LevelWarn},
		{http.StatusBadGateway, LevelWarn},
		{http.StatusServiceUnavailable, LevelWarn},
		{http.StatusGatewayTimeout, LevelWarn},
	}

	for _, tc := range tt {
		require.Equal(t, tc.level, StatusLevel(tc.status), "expected status %d to have level %d", tc.status, tc.level)
	}
}
