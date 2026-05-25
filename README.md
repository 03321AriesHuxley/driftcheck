# driftcheck

Detects configuration drift between running Docker containers and their source Compose definitions.

---

## Installation

```bash
go install github.com/yourusername/driftcheck@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/driftcheck.git
cd driftcheck
go build -o driftcheck .
```

---

## Usage

Point `driftcheck` at a Compose file and it will compare the defined configuration against your currently running containers.

```bash
# Check for drift using a docker-compose.yml in the current directory
driftcheck check

# Specify a custom Compose file
driftcheck check -f /path/to/docker-compose.yml

# Output results in JSON format
driftcheck check --output json
```

Example output:

```
[DRIFT] service "api": image mismatch (expected: myapp:1.2, running: myapp:1.1)
[DRIFT] service "worker": environment variable APP_ENV changed (expected: production, running: staging)
[OK]    service "db": no drift detected
```

Exit code `0` means no drift was found. Exit code `1` indicates drift was detected.

---

## How It Works

`driftcheck` reads your Compose file, queries the Docker daemon for running containers, and compares key properties including image tags, environment variables, port bindings, volume mounts, and resource limits.

---

## Contributing

Pull requests and issues are welcome. Please open an issue before submitting large changes.

---

## License

MIT © 2024 yourusername