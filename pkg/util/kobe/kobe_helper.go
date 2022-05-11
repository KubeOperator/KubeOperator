package kobe

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/KubeOperator/KubeOperator/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func NewKobeClient(host string, port int) *KobeClient {
	return &KobeClient{
		host: host,
		port: port,
	}
}

type KobeClient struct {
	host string
	port int
}

func (c *KobeClient) CreateProject(name string, source string) (*api.Project, error) {
	conn, err := c.createConnection()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := api.NewKobeApiClient(conn)
	request := api.CreateProjectRequest{
		Name:   name,
		Source: source,
	}
	resp, err := client.CreateProject(context.Background(), &request)
	if err != nil {
		return nil, err
	}
	return resp.Item, nil

}

func (c KobeClient) ListProject() ([]*api.Project, error) {
	conn, err := c.createConnection()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := api.NewKobeApiClient(conn)
	request := api.ListProjectRequest{}
	resp, err := client.ListProject(context.Background(), &request)
	if err != nil {
		return nil, err
	}
	return resp.Items, nil
}

func (c KobeClient) RunPlaybook(project, playbook, tag string, inventory *api.Inventory) (*api.KobeResult, error) {
	conn, err := c.createConnection()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := api.NewKobeApiClient(conn)
	request := &api.RunPlaybookRequest{
		Project:   project,
		Playbook:  playbook,
		Inventory: inventory,
		Tag:       tag,
	}
	req, err := client.RunPlaybook(context.Background(), request)
	if err != nil {
		return nil, err
	}
	return req.Result, nil
}

func (c KobeClient) RunAdhoc(pattern, module, param string, inventory *api.Inventory) (*api.KobeResult, error) {
	conn, err := c.createConnection()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := api.NewKobeApiClient(conn)
	request := &api.RunAdhocRequest{
		Inventory: inventory,
		Module:    module,
		Param:     param,
		Pattern:   pattern,
	}
	req, err := client.RunAdhoc(context.Background(), request)
	if err != nil {
		return nil, err
	}
	return req.Result, nil
}

func (c *KobeClient) WatchRun(taskId string, writer io.Writer) error {
	conn, err := c.createConnection()
	if err != nil {
		return err
	}
	defer conn.Close()
	client := api.NewKobeApiClient(conn)
	req := &api.WatchRequest{
		TaskId: taskId,
	}
	server, err := client.WatchResult(context.Background(), req)
	if err != nil {
		return err
	}
	for {
		msg, err := server.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		_, err = writer.Write(msg.Stream)
		if err != nil {
			break
		}
	}
	return nil
}

func (c *KobeClient) GetResult(taskId string) (*api.KobeResult, error) {
	conn, err := c.createConnection()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := api.NewKobeApiClient(conn)
	request := api.GetResultRequest{
		TaskId: taskId,
	}
	resp, err := client.GetResult(context.Background(), &request)
	if err != nil {
		return nil, err
	}
	return resp.Item, nil
}

func (c *KobeClient) ListResult() ([]*api.KobeResult, error) {
	conn, err := c.createConnection()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	client := api.NewKobeApiClient(conn)
	request := api.ListResultRequest{}
	resp, err := client.ListResult(context.Background(), &request)
	if err != nil {
		return nil, err
	}
	return resp.Items, nil
}

func (c *KobeClient) createConnection() (*grpc.ClientConn, error) {
	address := fmt.Sprintf("%s:%d", c.host, c.port)
	cre, err := newClientTLSFromFile("/var/kobe/conf/server.p", "kubeoperator.io")
	if err != nil {
		fmt.Printf("credentials.NewClientTLSFromFile err: %v", err)
	}
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(cre), grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(100*1024*1024)))
	if err != nil {
		fmt.Printf("grpc.Dial err: %v", err)
	}

	return conn, nil
}

func newClientTLSFromFile(certFile, serverNameOverride string) (credentials.TransportCredentials, error) {
	b, err := ioutil.ReadFile(certFile)
	if err != nil {
		return nil, err
	}

	cc, err := base64.StdEncoding.DecodeString(string(b))
	if err != nil {
		return nil, err
	}

	cp := x509.NewCertPool()
	if !cp.AppendCertsFromPEM(cc) {
		return nil, fmt.Errorf("credentials: failed to append certificates")
	}
	return credentials.NewTLS(&tls.Config{ServerName: serverNameOverride, RootCAs: cp}), nil
}
