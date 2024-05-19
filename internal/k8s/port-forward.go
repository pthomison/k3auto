package k8s

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	apimachinerytypes "k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/transport/spdy"
	"k8s.io/kubectl/pkg/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func PortForward(ctx context.Context, kcfg *rest.Config, name apimachinerytypes.NamespacedName, port int) (chan struct{}, error) {
	kcfg.GroupVersion = &schema.GroupVersion{Group: "api", Version: "v1"}
	kcfg.APIPath = "/"

	kcfg.NegotiatedSerializer = serializer.WithoutConversionCodecFactory{CodecFactory: scheme.Codecs}

	rc, err := rest.RESTClientFor(kcfg)
	if err != nil {
		return nil, err
	}

	req := rc.Post().
		Resource("pods").
		Namespace(name.Namespace).
		Name(name.Name).
		SubResource("portforward")

	transport, upgrader, err := spdy.RoundTripperFor(kcfg)
	if err != nil {
		return nil, err
	}

	stopChan := make(chan struct{})
	readyChan := make(chan struct{})

	host := kcfg.Host
	match, err := regexp.MatchString(`https://0\.0\.0\.0:\d{1,4}`, host)
	if err != nil {
		return nil, err
	}
	if match {
		host = strings.Replace(host, "https://0.0.0.0", "127.0.0.1", -1)
	} else {
		host = strings.Replace(host, "https://", "", -1)
	}

	u := req.URL()
	u.Host = host

	dialer := spdy.NewDialer(upgrader, &http.Client{Transport: transport}, "POST", u)

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

func PortForwardService(ctx context.Context, k8sC client.Client, kcfg *rest.Config, name apimachinerytypes.NamespacedName, port int) (chan struct{}, error) {
	svc := &corev1.Service{}
	err := k8sC.Get(ctx, name, svc)
	if err != nil {
		return nil, err
	}

	selector, err := metav1.LabelSelectorAsSelector(&metav1.LabelSelector{
		MatchLabels: svc.Spec.Selector,
	})
	if err != nil {
		return nil, err
	}

	pods := &corev1.PodList{}
	err = k8sC.List(ctx, pods, &client.ListOptions{
		LabelSelector: selector,
	})
	if err != nil {
		return nil, err
	}

	if len(pods.Items) < 1 {
		return nil, errors.New("no pods discovered for service")
	}

	podName := apimachinerytypes.NamespacedName{
		Name:      pods.Items[0].Name,
		Namespace: pods.Items[0].Namespace,
	}

	return PortForward(ctx, kcfg, podName, port)
}

func PortForwardDeployment(ctx context.Context, k8sC client.Client, kcfg *rest.Config, name apimachinerytypes.NamespacedName, port int) (chan struct{}, error) {
	dep := &appsv1.Deployment{}
	err := k8sC.Get(ctx, name, dep)
	if err != nil {
		return nil, err
	}

	selector, err := metav1.LabelSelectorAsSelector(dep.Spec.Selector)
	if err != nil {
		return nil, err
	}

	pods := &corev1.PodList{}
	err = k8sC.List(ctx, pods, &client.ListOptions{
		LabelSelector: selector,
	})
	if err != nil {
		return nil, err
	}

	if len(pods.Items) < 1 {
		return nil, errors.New("no pods discovered for deployment")
	}

	podName := apimachinerytypes.NamespacedName{
		Name:      pods.Items[0].Name,
		Namespace: pods.Items[0].Namespace,
	}

	return PortForward(ctx, kcfg, podName, port)
}
