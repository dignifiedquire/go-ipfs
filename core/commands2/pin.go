package commands

import (
	"fmt"

	cmds "github.com/jbenet/go-ipfs/commands"
	"github.com/jbenet/go-ipfs/core"
	"github.com/jbenet/go-ipfs/core/commands2/internal"
	"github.com/jbenet/go-ipfs/merkledag"
)

var pinCmd = &cmds.Command{
	Subcommands: map[string]*cmds.Command{
		"add": addPinCmd,
		"rm":  rmPinCmd,
	},
}

var addPinCmd = &cmds.Command{
	Options: []cmds.Option{
		cmds.Option{[]string{"recursive", "r"}, cmds.Bool},
		cmds.Option{[]string{"depth", "d"}, cmds.Uint},
	},
	Arguments: []cmds.Argument{
		cmds.Argument{"object", cmds.ArgString, true, true},
	},
	Run: func(res cmds.Response, req cmds.Request) {
		n := req.Context().Node

		// set recursive flag
		opt, _ := req.Option("recursive")
		recursive, _ := opt.(bool) // false if cast fails.

		paths, err := internal.ToStrings(req.Arguments())
		if err != nil {
			res.SetError(err, cmds.ErrNormal)
			return
		}

		_, err = pin(n, paths, recursive)
		if err != nil {
			res.SetError(err, cmds.ErrNormal)
		}

		// TODO: create some output to show what got pinned
	},
}

var rmPinCmd = &cmds.Command{
	Options: []cmds.Option{
		cmds.Option{[]string{"recursive", "r"}, cmds.Bool},
	},
	Arguments: []cmds.Argument{
		cmds.Argument{"object", cmds.ArgString, true, true},
	},
	Run: func(res cmds.Response, req cmds.Request) {
		n := req.Context().Node

		// set recursive flag
		opt, _ := req.Option("recursive")
		recursive, _ := opt.(bool) // false if cast fails.

		paths, err := internal.ToStrings(req.Arguments())
		if err != nil {
			res.SetError(err, cmds.ErrNormal)
			return
		}

		_, err = unpin(n, paths, recursive)
		if err != nil {
			res.SetError(err, cmds.ErrNormal)
		}

		// TODO: create some output to show what got unpinned
	},
}

func pin(n *core.IpfsNode, paths []string, recursive bool) ([]*merkledag.Node, error) {

	dagnodes := make([]*merkledag.Node, 0)
	for _, path := range paths {
		dagnode, err := n.Resolver.ResolvePath(path)
		if err != nil {
			return nil, fmt.Errorf("pin error: %v", err)
		}
		dagnodes = append(dagnodes, dagnode)
	}

	for _, dagnode := range dagnodes {
		err := n.Pinning.Pin(dagnode, recursive)
		if err != nil {
			return nil, fmt.Errorf("pin: %v", err)
		}
	}

	err := n.Pinning.Flush()
	if err != nil {
		return nil, err
	}

	return dagnodes, nil
}

func unpin(n *core.IpfsNode, paths []string, recursive bool) ([]*merkledag.Node, error) {

	dagnodes := make([]*merkledag.Node, 0)
	for _, path := range paths {
		dagnode, err := n.Resolver.ResolvePath(path)
		if err != nil {
			return nil, err
		}
		dagnodes = append(dagnodes, dagnode)
	}

	for _, dagnode := range dagnodes {
		k, _ := dagnode.Key()
		err := n.Pinning.Unpin(k, recursive)
		if err != nil {
			return nil, err
		}
	}

	err := n.Pinning.Flush()
	if err != nil {
		return nil, err
	}
	return dagnodes, nil
}
