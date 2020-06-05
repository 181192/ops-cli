package kubernetes

import (
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/signal"

	v1 "k8s.io/api/core/v1"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/transport/spdy"
)

// ExecClient is an interface for remote execution
type ExecClient interface {
	PodsForSelector(namespace, labelSelector string) (*v1.PodList, error)
	BuildPortForwarder(podName string, ns string, localPort int, podPort int) (*PortForward, error)
}

// PortForward gathers port forwarding results
type PortForward struct {
	Forwarder    *portforward.PortForwarder
	LocalPort    int
	StopChannel  chan struct{}
	ReadyChannel <-chan struct{}
}

// PodsForSelector get pods via label selector
func (client *Client) PodsForSelector(namespace, labelSelector string) (*v1.PodList, error) {
	podGet := client.Get().Resource("pods").Namespace(namespace).Param("labelSelector", labelSelector)
	obj, err := podGet.Do(context.Background()).Get()
	if err != nil {
		return nil, fmt.Errorf("failed retrieving pod: %v", err)
	}
	return obj.(*v1.PodList), nil
}

// BuildPortForwarder sets up port forwarding.
func (client *Client) BuildPortForwarder(podName string, ns string, localPort int, podPort int) (*PortForward, error) {
	var err error
	if localPort == 0 {
		localPort, err = availablePort()
		if err != nil {
			return nil, fmt.Errorf("failure allocating port: %v", err)
		}
	}
	req := client.Post().Resource("pods").Namespace(ns).Name(podName).SubResource("portforward")

	transport, upgrader, err := spdy.RoundTripperFor(client.Config)
	if err != nil {
		return nil, fmt.Errorf("failure creating roundtripper: %v", err)
	}

	dialer := spdy.NewDialer(upgrader, &http.Client{Transport: transport}, "POST", req.URL())

	stop := make(chan struct{})
	ready := make(chan struct{})
	fw, err := portforward.New(dialer, []string{fmt.Sprintf("%d:%d", localPort, podPort)}, stop, ready, ioutil.Discard, os.Stderr)
	if err != nil {
		return nil, fmt.Errorf("failed establishing port-forward: %v", err)
	}

	// Run the same check as k8s.io/kubectl/pkg/cmd/portforward/portforward.go
	// so that we will fail early if there is a problem contacting API server.
	podGet := client.Get().Resource("pods").Namespace(ns).Name(podName)
	obj, err := podGet.Do(context.Background()).Get()
	if err != nil {
		return nil, fmt.Errorf("failed retrieving pod: %v", err)
	}
	pod, ok := obj.(*v1.Pod)
	if !ok {
		return nil, fmt.Errorf("failed getting pod: %v", err)
	}
	if pod.Status.Phase != v1.PodRunning {
		return nil, fmt.Errorf("pod is not running. Status=%v", pod.Status.Phase)
	}

	return &PortForward{
		Forwarder:    fw,
		ReadyChannel: ready,
		StopChannel:  stop,
		LocalPort:    localPort,
	}, nil
}

func availablePort() (int, error) {
	addr, err := net.ResolveTCPAddr("tcp", ":0")
	if err != nil {
		return 0, err
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return 0, err
	}
	port := l.Addr().(*net.TCPAddr).Port
	return port, l.Close()
}

// RunPortForwarder run the portforwarder
func RunPortForwarder(fw *PortForward, readyFunc func(fw *PortForward) error) error {

	errCh := make(chan error, 1)
	go func() {
		errCh <- fw.Forwarder.ForwardPorts()
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	defer signal.Stop(signals)

	go func() {
		<-signals
		if fw.StopChannel != nil {
			close(fw.StopChannel)
		}
	}()

	select {
	case err := <-errCh:
		return fmt.Errorf("failure running port forward process: %v", err)
	case <-fw.ReadyChannel:
		err := readyFunc(fw)
		if err != nil {
			return err
		}

		// wait for interrupt (or connection close)
		<-fw.StopChannel
		return nil
	}
}
