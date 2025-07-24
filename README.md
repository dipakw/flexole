# flexole

**Flexole** is a lightweight, fast, and secure reverse proxy written in Go. It supports encrypted communication between local and remote servers, uses quantum-safe authentication, can be embedded into Go applications, and is designed to work in multi-user environments.

## Key Features

1. **Quantum-safe authentication**  
   Uses ML-KEM (Kyber) to provide post-quantum secure authentication, protecting against future quantum attacks.

2. **Multi-user support**  
   The server supports multiple users, each with their own access and configuration.

3. **Multiplexed communication**  
   Instead of opening a new connection for each request, Flexole uses a single connection to handle multiple requests efficiently.

4. **Configurable number of connections (pipes)**  
   Each service can use one or more dedicated connections (pipes). When multiple pipes are configured, requests are distributed using round-robin.

## Getting started
Download the appropriate binary file for your OS.

### Quick

```bash
flexole server --quick
```

Short
```
flexole s -q
```

### Expose local service
```bash
flexole client --quick --local=tcp/8080 --remote=tcp/80@server.com --id=100
```

Short
```bash
flexole c -q -l=tcp/8080 -r=tcp/80@server.com -i=100
```

-----------

### Using config files

#### Server

Generate a server config file, it generates a sample `server.yml`.

```bash
flexole generate --server-config
```

Or if you want a custom filename

```bash
flexole generate --server-config=myserver.yml
```

Start the server.

```bash
flexole server
```

Or to use a custom config file

```bash
flexole server --config=myserver.yml
```

#### Client

Generate a client config file, it generates a sample `client.yml`.

```bash
flexole generate --client-config
```

Or if you want a custom filename

```bash
flexole generate --client-config=myclient.yml
```

Start the client.

```bash
flexole client
```

Or to use a custom config file

```bash
flexole client --config=myclient.yml
```

## CLI usage

```
Usage:
  flexole <command> [options]

Commands:
  version, v     Show version
  server, s      Start the server
  client, c      Start the client
  generate, g    Generate sample config files
  help, h        Show this help message

Options (server):
  --config, -c    Path to server config file (default: server.yml)
  --quick, -q     Quick start with user key
  --dir, -d       Directory for server data (default: system temp dir)
  --user, -u      User ID (default: quick)
  --log, -o       Log levels: i=info, w=warn, e=error (default: iwe)
  --host, -h      Server host (default: 0.0.0.0)
  --port, -p      Server port (default: 8887)

Options (client):
  --config, -c     Path to client config file (default: client.yml)
  --quick, -q      Quick start with user key
  --user, -u       User ID (default: quick)
  --log, -o        Log levels: i=info, w=warn, e=error (default: iwe)
  --local, -l      Local address, e.g. tcp/localhost:8080
  --remote, -r     Remote, e.g. tcp/80@192.168.1.100
  --id             Service ID (0-65535)
  --encrypt, -e    Enable encryption (default: 1)

Options (generate):
  --client-config, -cc  Output client config file (default: client.yml)
  --server-config, -sc  Output server config file (default: server.yml)

Examples:
  flexole server --config=server.yml
  flexole client --quick=mykey --local=tcp/127.0.0.1:8080 --remote=tcp/80@192.168.1.100 --id=1
  flexole generate --client-config=client.yml --server-config=server.yml
  flexole help

Notes:
  - For log levels, combine letters (e.g., "iwe" for info, warn, error).
  - All options can use either --long or -short forms.
```