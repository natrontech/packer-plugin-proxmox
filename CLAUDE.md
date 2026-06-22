# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Build & install plugin locally (dev workflow)
make dev

# Run unit tests
make test

# Run acceptance tests (requires Proxmox environment)
PACKER_ACC=1 make testacc

# Run a single test
go test -v -run TestFoo ./builder/proxmox/common/

# Lint (via pre-commit)
pre-commit run --all-files

# Regenerate docs and HCL2 spec files
make generate

# Validate plugin structure
make plugin-check
```

## Architecture

This is a [Packer](https://www.packer.io/) plugin that builds Proxmox VM templates. It exposes two builders:
- `proxmox-iso` — creates a VM from scratch by booting an ISO
- `proxmox-clone` — clones an existing VM template

### Three-layer structure

```
iso/builder.go  clone/builder.go     ← thin wrappers; own only builder-specific config
        │               │
        └───────┬───────┘
        common/builder.go             ← shared multistep runner; owns Proxmox client lifecycle
                │
        common/step_*.go              ← generic pipeline steps
        clone/step_*.go               ← clone-specific steps
```

**`main.go`** registers both builders and is the plugin entry point.

**`common/config.go`** (~500 lines) defines the shared `Config` struct via composition: `common.PackerConfig`, `commonsteps.HTTPConfig`, `bootcommand.BootConfig`, `communicator.Config`, plus all Proxmox-specific options (disks, NICs, PCI passthrough, EFI, TPM, Cloud-Init, etc.).

**`common/step_start_vm.go`** (~821 lines) is the most complex file — it creates the VM, maps disk devices (IDE/SATA/SCSI/VirtIO), and generates the full Qemu config passed to the Proxmox API.

### VM creation strategy pattern

`stepStartVM` delegates the actual creation call to a `ProxmoxVMCreator` interface:
- `isoVMCreator` (in `iso/`) calls `config.Create(...)` on the Proxmox API
- `cloneVMCreator` (in `clone/`) calls `config.CloneVm(...)` on the Proxmox API

### Multistep pipeline

Both builders use a `multistep.Runner`. The state bag carries:
- `config` — build configuration
- `proxmoxClient` — Proxmox API client (`github.com/Telmate/proxmox-api-go`)
- `vmRef` — reference to the created/cloned VM
- `template_id` — set by `stepConvertToTemplate` at the end

Clone-specific steps run before the shared steps: `StepSshKeyPair` (generates ephemeral SSH keys for Cloud-Init injection) and `StepMapSourceDisks`.

### Testability

Steps that touch the Proxmox API accept narrow interfaces (`vmStarter`, `templateConverter`, `startedVMCleaner`) defined in `common/step_start_vm.go` and `common/step_convert_to_template.go`. Tests implement these interfaces with fakes.

## Module path

`github.com/natrontech/packer-plugin-proxmox` — this is a fork of `hashicorp/packer-plugin-proxmox`.

## Proxmox API reference

When adding new VM/disk/NIC features, the authoritative source for what fields are available is the `proxmox-api-go` library (the Go client, not the Proxmox HTTP API docs):

```
$(go env GOPATH)/pkg/mod/github.com/!telmate/proxmox-api-go@<version>/proxmox/
```

Key files for disk features:
- `config_qemu_disk.go` — shared `qemuDisk` struct and serialisation logic
- `config_qemu_disk_scsi.go` / `_ide.go` / `_sata.go` / `_virtio.go` — per-bus disk structs

Check the struct fields and inline comments (e.g. `// Only set for scsi,virtio`) to understand which options are valid per disk type before implementing config fields.

After adding new config fields, always run `make generate` to regenerate the HCL2 spec and docs files.

## Testing

- Add or update tests whenever a config field or feature is added or changed.
- Config validation tests go in `common/config_test.go` (use `TestDiskConfig` as a model for disk options).
- Disk wiring tests go in `common/step_start_vm_test.go` in `TestGenerateProxmoxDisks`.
