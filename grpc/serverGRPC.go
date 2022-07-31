package grpc

import (
	"context"
	"errors"
	"net"
	"os"

	sddaemon "github.com/coreos/go-systemd/v22/daemon"
	"github.com/seoyhaein/go-connections/sockets"
	"github.com/seoyhaein/utils"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

// TODO 살펴 보자. tls 도 지원해줘야 한다. 책 참고.
// https://hamait.tistory.com/931 살펴보고, service 등록하고, sddaemon 기능 파악하자.
// errgroup 기 능파악하자.

func ServeGRPC(cfg GRPCConfig, server *grpc.Server, errCh chan error) error {
	addrs := cfg.Address
	//if len(addrs) == 0 {
	if utils.IsEmptyString {
		return errors.New("--addr cannot be empty")
	}

	/*	tlsConfig, err := serverCredentials(cfg.TLS)
		if err != nil {
			return err
		}*/

	eg, _ := errgroup.WithContext(context.Background())
	listeners := make([]net.Listener, 0, len(addrs))
	for _, addr := range addrs {
		l, err := sockets.NewUnixSocketA(addr, *cfg.UID, *cfg.GID)
		if err != nil {
			for _, l := range listeners {
				l.Close()
			}
			return err
		}
		listeners = append(listeners, l)
	}

	if os.Getenv("NOTIFY_SOCKET") != "" {
		notified, notifyErr := sddaemon.SdNotify(false, sddaemon.SdNotifyReady)
		logrus.Debugf("SdNotifyReady notified=%v, err=%v", notified, notifyErr)
	}
	for _, l := range listeners {
		func(l net.Listener) {
			eg.Go(func() error {
				defer l.Close()
				logrus.Infof("running server on %s", l.Addr())
				return server.Serve(l)
			})
		}(l)
	}
	go func() {
		errCh <- eg.Wait()
	}()
	return nil
}

func ServeGRPCTls(cfg GRPCConfig, server *grpc.Server, errCh chan error) error {
	return nil
}
