package k8s

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/transport/spdy"
	"k8s.io/kubectl/pkg/scheme"
	ctlrconfig "sigs.k8s.io/controller-runtime/pkg/client/config"
)

func PortForward(ctx context.Context, name string, namespace string, port int) (chan struct{}, error) {
	kcfg, err := ctlrconfig.GetConfig()
	if err != nil {
		return nil, err
	}

	kcfg.GroupVersion = &schema.GroupVersion{Group: "api", Version: "v1"}
	kcfg.APIPath = "/"

	kcfg.NegotiatedSerializer = serializer.WithoutConversionCodecFactory{CodecFactory: scheme.Codecs}

	rc, err := rest.RESTClientFor(kcfg)
	if err != nil {
		return nil, err
	}

	req := rc.Post().
		Resource("pods").
		Namespace(namespace).
		Name(name).
		SubResource("portforward")

	transport, upgrader, err := spdy.RoundTripperFor(kcfg)
	if err != nil {
		return nil, err
	}

	stopChan := make(chan struct{})
	readyChan := make(chan struct{})

	u := req.URL()
	u.Host = "127.0.0.1:6443"

	dialer := spdy.NewDialer(upgrader, &http.Client{Transport: transport}, "POST", u)
	// if cmdutil.PortForwardWebsockets.IsEnabled() {
	// 	tunnelingDialer, err := portforward.NewSPDYOverWebsocketDialer(req.URL(), kcfg)
	// 	assert.Nil(t, err)

	// 	// First attempt tunneling (websocket) dialer, then fallback to spdy dialer.
	// 	dialer = portforward.NewFallbackDialer(tunnelingDialer, dialer, httpstream.IsUpgradeFailure)
	// }

	// url := strings.Replace(req.URL().String(), "0.0.0.0", , -1)

	// spew.Dump(url)

	fw, err := portforward.New(dialer, []string{fmt.Sprintf("%v:%v", port, port)}, stopChan, readyChan, os.Stdout, os.Stderr)
	if err != nil {
		return nil, err
	}

	go func() {
		err = fw.ForwardPorts()
	}()

	<-readyChan

	return stopChan, nil
}
