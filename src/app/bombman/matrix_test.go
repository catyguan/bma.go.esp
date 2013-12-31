package bombman

import (
	"logger"
	"testing"
	"time"
)

func initMatrix(m *Matrix) error {
	m.DoInit(4, 11, 11)
	logger.Info("TEST", "\n%s", m.View())
	return nil
}

func TestMatrix(t *testing.T) {
	m := NewMatrix(32)
	m.Run(initMatrix)

	m.AskClose()
	time.Sleep(1 * time.Millisecond)
}
