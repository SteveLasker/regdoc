package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/containerd/containerd/remotes"
	"github.com/containerd/containerd/remotes/docker"
	"github.com/stevelasker/regdoc/pkg/regdoc"
	"github.com/urfave/cli/v2"
)

var (
	Version  string
	Revision string

	ErrorNoPipe        = errors.New("this command is intended to work with pipes (|)")
	ErrorMissingRefArg = errors.New("missing arguments: [ref]")
)

var (
	usernameFlag = &cli.StringFlag{
		Name:    "username",
		Aliases: []string{"u"},
		Usage:   "username for generic remote access",
	}
	passwordFlag = &cli.StringFlag{
		Name:    "password",
		Aliases: []string{"p"},
		Usage:   "password for generic remote access",
	}
)

func main() {
	app := cli.NewApp()
	app.Name = "regdoc (regdoc)"
	app.Version = fmt.Sprintf("%s (build %s)", Version, Revision)
	app.Usage = "Registry documentation, persisting a repositories documentation within an OCI registry"
	app.Commands = []*cli.Command{
		{
			Name:  "push",
			Usage: "upload a markdown file to registry from stdin",
			Flags: []cli.Flag{
				usernameFlag,
				passwordFlag,
			},
			Action: func(c *cli.Context) error {
				ref := c.Args().Get(0)
				if ref == "" {
					return ErrorMissingRefArg
				}
				content, err := getStdin()
				if err != nil {
					return err
				}
				return regdoc.Push(content, ref, getResolver(c))
			},
		},
		{
			Name:  "pull",
			Usage: "download a repositories markdown file from registry and print to stdout",
			Flags: []cli.Flag{
				usernameFlag,
				passwordFlag,
			},
			Action: func(c *cli.Context) error {
				ref := c.Args().Get(0)
				if ref == "" {
					return ErrorMissingRefArg
				}
				return regdoc.Pull(ref, getResolver(c))
			},
		},
	}
	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func getStdin() ([]byte, error) {
	info, err := os.Stdin.Stat()
	if err != nil {
		return nil, err
	}
	if info.Mode()&os.ModeCharDevice != 0 || info.Size() <= 0 {
		return nil, ErrorNoPipe
	}
	reader := bufio.NewReader(os.Stdin)
	var output []rune
	for {
		input, _, err := reader.ReadRune()
		if err != nil && err == io.EOF {
			break
		}
		output = append(output, input)
	}
	return []byte(string(output)), nil
}

func getResolver(c *cli.Context) remotes.Resolver {
	opts := docker.ResolverOptions{}
	username := c.String(usernameFlag.Name)
	password := c.String(passwordFlag.Name)
	if username != "" || password != "" {
		opts.Credentials = func(hostName string) (string, string, error) {
			return username, password, nil
		}
	}
	return docker.NewResolver(opts)
}
