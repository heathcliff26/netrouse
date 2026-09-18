[![CI](https://github.com/heathcliff26/netrouse/actions/workflows/ci.yaml/badge.svg?event=push)](https://github.com/heathcliff26/netrouse/actions/workflows/ci.yaml)
[![Coverage Status](https://coveralls.io/repos/github/heathcliff26/netrouse/badge.svg)](https://coveralls.io/github/heathcliff26/netrouse)
[![Editorconfig Check](https://github.com/heathcliff26/netrouse/actions/workflows/editorconfig-check.yaml/badge.svg?event=push)](https://github.com/heathcliff26/netrouse/actions/workflows/editorconfig-check.yaml)
[![Generate go test cover report](https://github.com/heathcliff26/netrouse/actions/workflows/go-testcover-report.yaml/badge.svg)](https://github.com/heathcliff26/netrouse/actions/workflows/go-testcover-report.yaml)

# NetRouse

This is a simple utility for sending Wake-On-Lan magic packet to clients in the local network.
It can be used directly via the cli, or remotely via a web interface.

## Table of Contents

- [NetRouse](#netrouse)
  - [Table of Contents](#table-of-contents)
  - [Usage](#usage)
    - [CLI Args](#cli-args)
    - [Using the image](#using-the-image)
      - [Permissions for ping functionality](#permissions-for-ping-functionality)
    - [Image location](#image-location)
    - [Tags](#tags)
  - [Configuration](#configuration)
  - [Development](#development)
  - [Credit](#credit)

## Usage

### CLI Args
```bash
$ netrouse help
NetRouse power on other devices on the network via Wake-on-Lan

Usage:
  netrouse [flags]
  netrouse [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  server      Serve a frontend via gui
  version     Print version information and exit
  wol         Send a magic packet to the given mac address

Flags:
  -h, --help   help for netrouse

Use "netrouse [command] --help" for more information about a command.
```

### Using the image

When using the container image, please note that the server needs to run with `--net host` to send the magic packets.
```bash
$ podman run -d --net host -v /path/to/config.yaml:/config/config.yaml ghcr.io/heathcliff26/netrouse:latest
```
If you want to use it with persistent data, run it with `-v netrouse-data:/data`. With the default configuration it will write data to `/data/hosts.yaml`.

The image can be run without configuration, it will simply use the default values.

#### Permissions for ping functionality

If you encounter `socket: permission denied` errors when checking if a host is online, you might need to run `sudo sysctl -w net.ipv4.ping_group_range="0 2147483647`.

### Image location

| Container Registry                                                                            | Image                             |
| --------------------------------------------------------------------------------------------- | --------------------------------- |
| [Github Container](https://github.com/users/heathcliff26/packages/container/package/netrouse) | `ghcr.io/heathcliff26/netrouse`   |
| [Docker Hub](https://hub.docker.com/r/heathcliff26/netrouse)                                  | `docker.io/heathcliff26/netrouse` |
| [Quay.io](https://quay.io/heathcliff26/netrouse)                                              | `quay.io/heathcliff26/netrouse`   |

### Tags

There are different flavors of the image:

| Tag(s)      | Description                                                 |
| ----------- | ----------------------------------------------------------- |
| **latest**  | Last released version of the image                          |
| **rolling** | Rolling update of the image, always build from main branch. |
| **vX.Y.Z**  | Released version of the image                               |

## Configuration

An example configuration with comments and default values for the server can be found [here](examples/config.yaml).

The default paths for the configuration file are:
    - standalone:   `/etc/netrouse/config.yaml`
    - in container: `/config/config.yaml`

## Development

When changing the frontend, use the original bootstrap file `static/bootstrap/bootstrap.css` and not the trimmed version `static/css/bootstrap.css`

## Credit

For CSS bootstrap is used: [Website](https://getbootstrap.com/) | [Github](https://github.com/twbs/bootstrap) | [Documentation](https://getbootstrap.com/docs/5.3/getting-started/introduction/)

To trim the bootstrap file purgecss is used: [Website](https://purgecss.com/) | [Github](https://github.com/FullHuman/purgecss) | [Documentation](https://purgecss.com/getting-started.html)

The favicon base is from bootstrap: [Power Button](https://icons.getbootstrap.com/icons/power/)

For generating the different icons: [RealFaviconGenerator](https://realfavicongenerator.net/)
