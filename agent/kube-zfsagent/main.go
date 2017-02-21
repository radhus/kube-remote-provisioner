package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/mistifyio/go-zfs"

	"github.com/radhus/kube-remote-provisioner/agent/kube-zfsagent/agent"
)

var (
	// gRPC options
	listenAddr = flag.String("listen", ":8080", "Host and port to listen on")

	// zfs options
	rootDataset = flag.String("dataset", "", "Root dataset to create datasets under")

	// nfs options
	nfsHost = flag.String("nfshost", "", "NFS hostname for Kubernetes to use")
)

func flagValidateOrTerminate() {
	fatal := func(err string) {
		fmt.Fprintln(os.Stderr, err)
		flag.Usage()
		os.Exit(1)
	}

	flag.Parse()

	// TODO: fix TLS setup
	if _, err := zfs.GetDataset(*rootDataset); err != nil {
		fatal(fmt.Sprintf("error getting dataset: %s", err))
	}
	// TODO: validate type of dataset?

	if *nfsHost == "" {
		fatal("nfshost is required")
	}
}

func main() {
	flagValidateOrTerminate()
	if err := agent.Run(
		*listenAddr,
		*rootDataset,
		*nfsHost,
	); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
