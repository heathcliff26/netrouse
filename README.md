[![CI](https://github.com/heathcliff26/netrouse/actions/workflows/ci.yaml/badge.svg?event=push)](https://github.com/heathcliff26/netrouse/actions/workflows/ci.yaml)
[![Coverage Status](https://coveralls.io/repos/github/heathcliff26/netrouse/badge.svg)](https://coveralls.io/github/heathcliff26/netrouse)
[![Editorconfig Check](https://github.com/heathcliff26/netrouse/actions/workflows/editorconfig-check.yaml/badge.svg?event=push)](https://github.com/heathcliff26/netrouse/actions/workflows/editorconfig-check.yaml)
[![Coverprofiles](https://github.com/heathcliff26/netrouse/actions/workflows/coverprofiles.yaml/badge.svg)](https://github.com/heathcliff26/netrouse/actions/workflows/coverprofiles.yaml)

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
  - [Installing the GUI](#installing-the-gui)
    - [Download binary](#download-binary)
      - [Uninstalling](#uninstalling)
    - [Fedora Copr](#fedora-copr)
    - [Android](#android)
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

## Installing the GUI

### Download binary

1. Download the [latest release](https://github.com/heathcliff26/netrouse/releases/latest)
2. Unpack the archive
3. Install the app for your user by running:
   - You can install it globally by running the script with `sudo`
```bash
./install.sh -i
```

#### Uninstalling

1. Switch to the folder where you have the installation script
2. Uninstall by running:
   - Run as `sudo` if you installed it globally
```bash
./install.sh -u
```
3. Delete the folder.


### Fedora Copr

The app is available as an rpm by using the fedora copr repository [heathcliff26/NetRouse](https://copr.fedorainfracloud.org/coprs/heathcliff26/NetRouse/).
1. Enable the copr repository
```bash
sudo dnf copr enable heathcliff26/NetRouse
```
2. Install the app
```bash
sudo dnf install netrouse
```

**Note:** The server/cli binary is also available on copr as `netroused`.

### Android

You can download the latest release from [here](https://github.com/heathcliff26/netrouse/releases/latest).
For smaller file size, choose the apk matching your phones cpu architecture.

Alternatively, you can install and automatically update the app by using [Obtainium](https://obtainium.imranr.dev/):

[![](images/badge_obtainium.png)](https://apps.obtainium.imranr.dev/redirect?r=obtainium://app/%7B%22id%22%3A%20%22io.github.heathcliff26.netrouse%22%2C%20%22url%22%3A%20%22https%3A%2F%2Fgithub.com%2Fheathcliff26%2Fnetrouse%22%2C%20%22author%22%3A%20%22Heathcliff%22%2C%20%22name%22%3A%20%22NetRouse%22%7D)

## Development

When changing the frontend, use the original bootstrap file `static/bootstrap/bootstrap.css` and not the trimmed version `static/css/bootstrap.css`

## Credit

For CSS bootstrap is used: [Website](https://getbootstrap.com/) | [Github](https://github.com/twbs/bootstrap) | [Documentation](https://getbootstrap.com/docs/5.3/getting-started/introduction/)

To trim the bootstrap file purgecss is used: [Website](https://purgecss.com/) | [Github](https://github.com/FullHuman/purgecss) | [Documentation](https://purgecss.com/getting-started.html)

The favicon base is from bootstrap: [Power Button](https://icons.getbootstrap.com/icons/power/)

For generating the different icons: [RealFaviconGenerator](https://realfavicongenerator.net/)
