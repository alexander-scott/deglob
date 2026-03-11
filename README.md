# deglob

**deglob** is a command-line tool that removes [glob patterns](https://bazel.build/reference/be/functions#glob) from [Bazel](https://bazel.build/) `BUILD` files, replacing them with explicit, individually listed file targets.

## Motivation

Bazel `glob()` patterns are convenient for pulling in many files at once, but they make build graphs harder to reason about, can hide unused dependencies, and create implicit coupling between files on disk and build targets. By replacing globs with explicit targets, you get:

- **Deterministic builds** — no surprises when a new file is added to a directory.
- **Clearer dependency graphs** — each file is an explicit node; it is obvious what each target depends on.
- **Easier refactoring** — removing or renaming a file requires a deliberate BUILD file change rather than a silent glob match change.

## How It Works

deglob scans every `BUILD` file in a workspace, finds targets that use `glob()` in their `srcs` or `hdrs` attributes, expands the glob pattern against the files that currently exist on disk, and rewrites the `BUILD` file so that:

1. The original target's glob attribute is replaced with a `deps` list pointing to new child targets.
2. One new `cc_library` target is emitted per matched file, each with a single, explicit file entry.

### Example

**Before** (`BUILD`):

```python
cc_library(
    name = "files_with_glob",
    hdrs = glob(["file*.cpp"]),
)

cc_binary(
    name = "main",
    srcs = ["main.cpp"],
    deps = [":files_with_glob"],
)
```

**After** (`BUILD`):

```python
cc_library(
    name = "files_with_glob",
    deps = [":files_with_glob_file_1", ":files_with_glob_file_2"],
)

cc_library(
    name = "files_with_glob_file_1",
    hdrs = ["file_1.cpp"],
)

cc_library(
    name = "files_with_glob_file_2",
    hdrs = ["file_2.cpp"],
)

cc_binary(
    name = "main",
    srcs = ["main.cpp"],
    deps = [":files_with_glob"],
)
```

## Prerequisites

| Requirement | Version |
|-------------|---------|
| [Go](https://go.dev/dl/) | 1.21+ |

The released binary has no runtime dependencies (built with `CGO_ENABLED=0`).

For local development you will also need:

- [Make](https://www.gnu.org/software/make/)
- [GoReleaser](https://goreleaser.com/install/)
- [golangci-lint](https://golangci-lint.run/welcome/install/)
- [misspell](https://github.com/client9/misspell)

## Installation

### Pre-built binaries

Download the latest release for your platform from the [Releases](https://github.com/alexander-scott/deglob/releases) page and place the binary somewhere on your `PATH`.

### Docker

```bash
docker pull ghcr.io/alexander-scott/deglob:latest
```

### Build from source

```bash
git clone https://github.com/alexander-scott/deglob.git
cd deglob
make build
# Binary is written to dist/
```

## Usage

```
deglob -workspace_path=<path> [-filter=<regex>]
```

| Flag | Required | Default | Description |
|------|----------|---------|-------------|
| `-workspace_path` | Yes | — | Path to the root of the Bazel workspace to process. |
| `-filter` | No | `.*BUILD$` | Regular expression used to select which files inside the workspace are processed. |

### Examples

Process all `BUILD` files in a workspace:

```bash
deglob -workspace_path=/path/to/my/workspace
```

Process only files named `BUILD.bazel`:

```bash
deglob -workspace_path=/path/to/my/workspace -filter=".*BUILD\.bazel$"
```

Run via Docker:

```bash
docker run --rm \
  -v /path/to/my/workspace:/workspace \
  ghcr.io/alexander-scott/deglob:latest \
  -workspace_path=/workspace
```

## Development

```bash
# Run the full build pipeline (format, lint, build, test)
make all

# Run tests only
make test

# Run the linter (with auto-fix)
make lint

# Run spell check on Markdown files
make spell

# Remove build artifacts
make clean
```

Test coverage reports are written to `coverage.html` after `make test`.
