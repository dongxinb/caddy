package caddyhttp

import (
	"unsafe"

	"go.uber.org/zap"
	"golang.org/x/sys/unix"
)

const (
	BrutalAvailable   = true
	TCP_BRUTAL_PARAMS = 23301
)

type TCPBrutalParams struct {
	Rate     uint64
	CwndGain uint32
}

//go:linkname setsockopt syscall.setsockopt
func setsockopt(s int, level int, name int, val unsafe.Pointer, vallen uintptr) (err error)

func EnableBrutalSockopt(logger *zap.Logger, fd uintptr) error {
	err := unix.SetsockoptString(int(fd), unix.IPPROTO_TCP, 0xd, "brutal")
	if err != nil {
		logger.Error("enable brutal failed", zap.Error(err))
		return err
	}
	logger.Info("enable brutal success")
	return nil
}

func SetBrutalSockopt(logger *zap.Logger, fd uintptr, brutal_speed int) error {
	params := TCPBrutalParams{
		Rate:     uint64(125000) * uint64(brutal_speed),
		CwndGain: 20, // hysteria2 default
	}
	err := setsockopt(int(fd), unix.IPPROTO_TCP, TCP_BRUTAL_PARAMS, unsafe.Pointer(&params), unsafe.Sizeof(params))
	if err != nil {
		logger.Error("brutal set speed failed", zap.Error(err))
		return err
	}
	return nil
}
