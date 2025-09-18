# oci-cluster-token

A token generator for accessing Oracle Cloud Infrastructure (OCI) Kubernetes Engine (OKE) clusters in `~/.kube/config`. This tool replaces the `oci` command to accelerate `kubectl` and `k9s` commands by generating tokens more efficiently.

## Overview

`oci-cluster-token` is a lightweight utility designed to streamline access to Oracle Cloud's Kubernetes OKE clusters. Instead of relying on the `oci` command-line tool to fetch tokens, this tool generates them directly, improving the speed of Kubernetes CLI operations like `kubectl` and `k9s`.

## Features

- Generates tokens for OKE cluster access.
- Replaces the `oci` command for faster `~/.kube/config` integration.
- Speeds up `kubectl` and `k9s` command execution.

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/mattn/oci-cluster-token.git
   ```
2. Build the tool:
   ```bash
   cd oci-cluster-token
   go build
   ```
3. Move the binary to your PATH:
   ```bash
   mv oci-cluster-token /usr/local/bin/
   ```

## Usage

1. Configure your `~/.kube/config` to use `oci-cluster-token` instead of the `oci` command for token generation.
2. Update the `exec` section in your `~/.kube/config` to reference `oci-cluster-token`:
   ```yaml
   - command: oci-cluster-token
     apiVersion: client.authentication.k8s.io/v1beta1
     args:
       - --cluster-id
       - <your-cluster-id>
   ```
3. Run `kubectl` or `k9s` as usual, and the tool will handle token generation.

## Installation

- Access to an Oracle Cloud Infrastructure account and an OKE cluster.
- Properly configured `~/.oci/config` with your OCI credentials.

```
go install github.com/mattn/oci-cluster-token@latest
```

## License

MIT

## Author

Yasuhiro Matsumoto (a.k.a. mattn)
