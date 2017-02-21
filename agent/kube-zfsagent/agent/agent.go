package agent

import (
	"errors"
	"fmt"
	"log"
	"net"
	"path"

	"github.com/mistifyio/go-zfs"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/radhus/kube-remote-provisioner/api"
)

type agent struct {
	rootDataset string
	nfsHost     string
}

var _ api.AgentServiceServer = (*agent)(nil)

func Run(
	listenAddr string,
	rootDataset string,
	nfsHost string,
) error {
	a := &agent{
		rootDataset: rootDataset,
		nfsHost:     nfsHost,
	}

	lis, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatal("Listen failed:", err)
	}

	s := grpc.NewServer()
	api.RegisterAgentServiceServer(s, a)

	log.Println("Listening for gRPC on", lis.Addr().String())

	return s.Serve(lis)
}

func (a *agent) getProvisionResponse(name string, readOnly bool) *api.ProvisionResponse {
	return &api.ProvisionResponse{
		Source: &api.Source{
			Type: &api.Source_Nfs{
				Nfs: &api.NFSVolumeSource{
					Server:   a.nfsHost,
					Path:     path.Join("/", name),
					ReadOnly: readOnly,
				},
			},
		},
	}
}

func (a *agent) Provision(ctx context.Context, req *api.ProvisionRequest) (*api.ProvisionResponse, error) {
	readOnly := false
	props := map[string]string{}

	if pvc := req.GetPvc(); pvc != nil {
		if spec := pvc.GetSpec(); spec != nil {
			// TODO: accessModes
			if resources := spec.GetResources(); resources != nil {
				if resources.Requests > 0 {
					// TODO: check free space
					// TODO: set reservation?
				}
				if resources.Limits > 0 {
					props["quota"] = fmt.Sprintf("%d", resources.Limits)
				}
			}
		}
	}

	// TODO: set other stuff as user properties

	// TODO: validate and sanitize input Name
	name := path.Join(a.rootDataset, req.Name)

	dataset, err := zfs.CreateFilesystem(name, props)
	if err != nil {
		log.Printf(
			"Error creating filesystem, name: %s, props: %v, error: %s",
			name, props, err,
		)
		return nil, err
	}

	log.Println("Created dataset:", name)

	response := a.getProvisionResponse(name, readOnly)
	response.Capacity = dataset.Avail

	return response, nil
}

func (a *agent) Delete(_ context.Context, req *api.DeleteRequest) (*api.DeleteResponse, error) {
	return nil, errors.New("not implemented yet")
}
