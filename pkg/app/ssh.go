// Copyright 2022 The envd Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package app

import (
	"github.com/cockroachdb/errors"
	"github.com/urfave/cli/v2"

	"github.com/tensorchord/envd/pkg/home"
	"github.com/tensorchord/envd/pkg/ssh"
	sshconfig "github.com/tensorchord/envd/pkg/ssh/config"
)

var CommandSSH = &cli.Command{
	Name:     "ssh",
	Category: CategoryBasic,
	Hidden:   true,
	Usage:    "TestK8s",
	Flags: []cli.Flag{
		&cli.PathFlag{
			Name:    "private-key",
			Usage:   "Path to the private key",
			Aliases: []string{"k"},
			Value:   sshconfig.GetPrivateKeyOrPanic(),
			Hidden:  true,
		},
		&cli.PathFlag{
			Name:    "public-key",
			Usage:   "Path to the public key",
			Aliases: []string{"pubk"},
			Value:   sshconfig.GetPublicKeyOrPanic(),
			Hidden:  true,
		},
	},
	Action: sshc,
}

func sshc(clicontext *cli.Context) error {
	ac, err := home.GetManager().AuthGetCurrent()
	if err != nil {
		return err
	}
	it := ac.IdentityToken

	opt := ssh.DefaultOptions()
	opt.User = it
	opt.PrivateKeyPath = clicontext.Path("private-key")
	opt.Port = 2222
	opt.AgentForwarding = false
	sshClient, err := ssh.NewClient(opt)
	if err != nil {
		return errors.Wrap(err, "failed to create the ssh client")
	}
	if err := sshClient.Attach(); err != nil {
		return errors.Wrap(err, "failed to attach to the container")
	}
	return nil
}