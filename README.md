# terracost

[![Actions Status](https://github.com/kainosnoema/terracost-cli/workflows/Test/badge.svg)](https://github.com/kainosnoema/terracost-cli/actions)

AWS cost estimation for Terraform projects using a custom API hosted at terracost.io. Does not read or upload any Terraform state, variables, or outputs.

## Installation

### Homebrew

```console
$ brew install kainosnoema/tap/terracost
```

### Usage

In Terraform project directory:
```
$ terracost estimate [plan file]
```
