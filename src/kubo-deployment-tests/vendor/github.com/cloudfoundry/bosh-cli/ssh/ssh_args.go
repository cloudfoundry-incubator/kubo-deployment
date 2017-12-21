package ssh

import (
	"fmt"
	"strings"

	boshsys "github.com/cloudfoundry/bosh-utils/system"

	boshdir "github.com/cloudfoundry/bosh-cli/director"
)

type SSHArgs struct {
	ConnOpts ConnectionOpts
	Result   boshdir.SSHResult

	ForceTTY bool

	PrivKeyFile    boshsys.File
	KnownHostsFile boshsys.File
}

func (a SSHArgs) LoginForHost(host boshdir.Host) []string {
	return []string{host.Host, "-l", host.Username}
}

func (a SSHArgs) OptsForHost(host boshdir.Host) []string {
	// Options are used for both ssh and scp
	cmdOpts := []string{}

	if a.ForceTTY {
		cmdOpts = append(cmdOpts, "-tt")
	}

	cmdOpts = append(cmdOpts, []string{
		"-o", "ServerAliveInterval=30",
		"-o", "ForwardAgent=no",
		"-o", "PasswordAuthentication=no",
		"-o", "IdentitiesOnly=yes",
		"-o", "IdentityFile=" + a.PrivKeyFile.Name(),
		"-o", "StrictHostKeyChecking=yes",
		"-o", "UserKnownHostsFile=" + a.KnownHostsFile.Name(),
	}...)

	gwUsername, gwHost, gwPrivKeyPath := a.gwOpts()

	if len(a.ConnOpts.SOCKS5Proxy) > 0 {
		proxyOpt := fmt.Sprintf(
			"ProxyCommand=nc -x %s %%h %%p",
			strings.TrimPrefix(a.ConnOpts.SOCKS5Proxy, "socks5://"),
		)

		cmdOpts = append(cmdOpts, "-o", proxyOpt)

	} else if len(gwHost) > 0 {
		gwCmdOpts := []string{
			"-o", "ServerAliveInterval=30",
			"-o", "ForwardAgent=no",
			"-o", "ClearAllForwardings=yes",
			// Strict host key checking for a gateway is not necessary
			// since ProxyCommand is only used for forwarding TCP and
			// agent forwarding is disabled
			"-o", "StrictHostKeyChecking=no",
			"-o", "UserKnownHostsFile=/dev/null",
		}

		if len(gwPrivKeyPath) > 0 {
			gwCmdOpts = append(
				gwCmdOpts,
				"-o", "PasswordAuthentication=no",
				"-o", "IdentitiesOnly=yes",
				"-o", "IdentityFile="+gwPrivKeyPath,
			)
		}

		// It appears that when using ssh -W, IPv6 address needs to be put in brackets
		// fixes: `Bad stdio forwarding specification 'fd7a:eeed:e696:...'`
		proxyHostPortTmpl := "%h:%p"
		if strings.Contains(host.Host, ":") {
			proxyHostPortTmpl = "[%h]:%p"
		}

		proxyOpt := fmt.Sprintf(
			// Always force TTY for gateway ssh
			"ProxyCommand=ssh -tt -W %s -l %s %s %s",
			proxyHostPortTmpl,
			gwUsername,
			gwHost,
			strings.Join(gwCmdOpts, " "),
		)

		cmdOpts = append(cmdOpts, "-o", proxyOpt)
	}

	cmdOpts = append(cmdOpts, a.ConnOpts.RawOpts...)

	return cmdOpts
}

func (a SSHArgs) gwOpts() (string, string, string) {
	if a.ConnOpts.GatewayDisable {
		return "", "", ""
	}

	// Take server provided gateway options
	username := a.Result.GatewayUsername
	host := a.Result.GatewayHost

	if len(a.ConnOpts.GatewayUsername) > 0 {
		username = a.ConnOpts.GatewayUsername
	}

	if len(a.ConnOpts.GatewayHost) > 0 {
		host = a.ConnOpts.GatewayHost
	}

	privKeyPath := a.ConnOpts.GatewayPrivateKeyPath

	return username, host, privKeyPath
}
