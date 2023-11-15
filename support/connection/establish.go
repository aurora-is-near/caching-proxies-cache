package connection

import (
	"fmt"

	"github.com/nats-io/jsm.go/natscontext"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

func Establish(context string, overrideServer string, overrideCreds string) (*nats.Conn, jetstream.JetStream, error) {
	var options = []natscontext.Option{}
	if overrideCreds != "" {
		options = append(options, natscontext.WithCreds(overrideCreds))
	}
	if overrideServer != "" {
		options = append(options, natscontext.WithServerURL(overrideServer))
	}

	ctx, err := natscontext.New(context, true, options...)
	if err != nil {
		return nil, nil, err
	}

	ns, err := ctx.Connect()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to nats: %w", err)
	}

	js, err := jetstream.New(ns)
	if err != nil {
		_ = ns.Drain()
		return nil, nil, fmt.Errorf("failed to connect to jetstream: %w", err)
	}

	return ns, js, err
}
