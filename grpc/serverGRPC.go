package grpc

import (
	"context"
	"errors"
	"net"
	"os"

	"github.com/seoyhaein/go-connections/sockets"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

func ServeGRPC(cfg config.GRPCConfig, server *grpc.Server, errCh chan error) error {
	addrs := cfg.Address
	if len(addrs) == 0 {
		return errors.New("--addr cannot be empty")
	}

	//tlsConfig, err := serverCredentials(cfg.TLS)
	if err != nil {
		return err
	}
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
