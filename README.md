# lima-compose

Compose for Lima VMs - A thin wrapper around limactl for declarative multi-VM management.

## Overview

lima-compose allows you to define and manage multiple Lima VMs using a YAML file. It acts as a thin wrapper around `limactl`, providing the following capabilities:

- Define multiple VM configurations in a single YAML file
- Execute batch operations (create/destroy/start/stop) on all defined VMs
- Display the exact `limactl` commands being executed
- Output VM IP addresses in `/etc/hosts` format

Unlike Docker Compose, lima-compose intentionally provides minimal orchestration features - no service dependencies, no restart policies, and no built-in provisioning. It simply translates YAML definitions into `limactl` commands.

## Features

- **Declarative YAML configuration** - Define all your VMs in a single file
- **Batch operations** - Create, start, stop, and destroy all VMs with one command
- **Transparent execution** - See exactly which `limactl` commands are being run
- **Zero learning curve** - If you know `limactl`, you already know lima-compose
- **Network information** - ExportVM IP addresses in `/etc/hosts` format

## Prerequisites

[Lima](https://lima-vm.io) installed and `limactl` available in PATH

## Installation

```bash
$ go install github.com/zinrai/lima-compose@latest
```

## YAML Configuration

The configuration file consists of an `instances` section where each VM is defined with:

- **template** (required): The Lima template to use (URL, path, or `template://` reference)
- **args** (optional): Command-line arguments passed directly to `limactl create`

### Minimal Example

```yaml
instances:
  my-vm:
    template: template://debian
    args: |
      --cpus 2
      --memory 2
```

### Full Example with Multiple VMs

```yaml
instances:
  consul-server-01:
    template: template://ubuntu-lts
    args: |
      --cpus 2
      --memory 4
      --disk 50
      --set '.env.NODE_TYPE="server"'
      
  consul-server-02:
    template: template://ubuntu-lts
    args: |
      --cpus 2
      --memory 4
      --disk 50
      --set '.env.NODE_TYPE="server"'
      
  consul-client-01:
    template: template://ubuntu-lts
    args: |
      --cpus 4
      --memory 8
      --disk 100
      --mount ~/projects:/projects:w
      --set '.env.NODE_TYPE="client"'
```

## Commands

### create

Creates all VMs defined in the YAML file.

```bash
$ lima-compose create [compose-file]
```

- Executes `limactl create` for each instance
- Shows the exact `limactl` command being run
- Fails if any VM cannot be created

### destroy

Destroys all VMs defined in the YAML file.

```bash
$ lima-compose destroy [compose-file]
```

### start

Starts all VMs defined in the YAML file.

```bash
$ lima-compose start [compose-file]
```

### stop

Stops all VMs defined in the YAML file.

```bash
$ lima-compose stop [compose-file]
```

### ips

Shows IPv4 addresses of all running VMs in `/etc/hosts` format, one line per interface.

```bash
$ lima-compose ips [compose-file]
```

Output example:
```
127.0.0.1       web-01-lo
192.168.5.15    web-01-eth0
127.0.0.1       db-01-lo
192.168.5.16    db-01-eth0
```

If no compose file is specified, `lima-compose.yaml` or `lima-compose.yml` is used by default.

### version

Shows the version, commit, and build date embedded at build time.

```bash
$ lima-compose version
```

## Usage Examples

### Basic Workflow

```bash
# Create VMs
lima-compose create lima-compose.yaml

# Check status using standard limactl
limactl list

# Get IP addresses
lima-compose ips lima-compose.yaml

# SSH into a specific VM using standard limactl
limactl shell web-01

# Stop all VMs
lima-compose stop lima-compose.yaml

# Start all VMs
lima-compose start lima-compose.yaml

# Destroy all VMs
lima-compose destroy lima-compose.yaml
```

### Distributing hosts Configuration

```bash
# Get IP addresses and save to file
lima-compose ips > /tmp/hosts

# Distribute to all VMs
for vm in web-01 db-01; do
  cat /tmp/hosts | limactl shell $vm "sudo tee -a /etc/hosts"
done
```

## Design Philosophy

lima-compose is intentionally designed as a **thin wrapper** around `limactl`:

- **Transparent**: Every `limactl` command executed is shown to the user
- **Minimal**: Only provides multi-VM management, nothing more
- **Non-invasive**: Does not modify or enhance Lima's functionality
- **Predictable**: YAML content directly maps to `limactl` commands

## Limitations

By design, lima-compose does NOT:

- Support partial operations (all VMs or none)
- Control startup order or dependencies
- Provide provisioning or configuration management
- Manage VM state beyond what `limactl` provides
- Hide or abstract Lima's behavior

For individual VM operations, use `limactl` directly.

## License

This project is licensed under the [MIT License](./LICENSE).
