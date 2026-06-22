# Testing

Always write or update tests when adding or changing a feature:

- **Config validation** (new fields, constraints): add cases to the relevant `TestXxx` in `common/config_test.go`. Use `TestDiskConfig` as a model for disk options.
- **Disk wiring** (new fields passed to the Proxmox API): add cases to `TestGenerateProxmoxDisks` in `common/step_start_vm_test.go` and update any existing expected structs that are affected.
- **New steps or step behaviour**: add or extend `Test*` functions in the corresponding `*_test.go` file.

Run `make test` before considering the task done.
