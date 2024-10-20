# terraform-lockfile-insights

A utility for surfacing information about the Terraform providers used in a project. It recursively crawls a supplied directory, and parses each Terraform lockfile using `tree-sitter` to provide a report showing which lockfiles use which versions of which providers.

Example output (using `--pretty`):
```json
{
  "registry.opentofu.org/cloudflare/cloudflare": {
    "versions": {
      "4.19.0": [
        "infrastructure/terraform/foo/.terraform.lock.hcl",
        "infrastructure/terraform/bar/.terraform.lock.hcl",
        "infrastructure/terraform/baz/.terraform.lock.hcl"
      ]
    }
  },
  "registry.opentofu.org/hashicorp/external": {
    "versions": {
      "2.3.3": [
        "infrastructure/terraform/foo/.terraform.lock.hcl",
        "infrastructure/terraform/bar/.terraform.lock.hcl",
        "infrastructure/terraform/baz/.terraform.lock.hcl"
      ]
    }
  }
}
```

## Usage

There are two ways to run `terraform-lockfile-insights`:

1. Via the Golang CLI: `go run main.go <flags> <repo_path>`
2. Via the Nix flake: `nix run . -- <flags> <repo_path>`

You can optionally supply a `--pretty` flag to pretty print the JSON
output.
