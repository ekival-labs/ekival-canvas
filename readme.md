# EKIVAL-Canvas: Backend Transaction Building APIs Testing for Ekival Server

## Introduction

EKIVAL-Canvas is a Golang backend service that provides APIs for building transactions and interacting with the Ekival server. It simplifies integration of Ekival functionality into your applications.

## Features

- **Transaction Building:** Generate and sign transactions for sending/receiving HNS, managing names, and voting on proposals.
- **Handshake Server Interaction:** Unify communication and data retrieval through a single API.
- **JSON Configuration:** Easily configure server settings and API keys with a dedicated JSON file.
- **Lightweight & Efficient:** Developed in Golang for fast and resource-friendly operation.

## Installation & Usage

1. **Dependencies:** Go 1.21+ (`go version`)
2. **Clone the Repository:** `git clone https://ekival-canvas.git`
3. **Change directory:** `cd ekival-canvas`
4. **Create build folder:** `mkdir build`
5. **Build the Binary:** `go build -o ./build/`
6. **Start the Server:** `./build/ekival-canvas -config config/config.json` (Listens on port **_42069_** by default)

## API Documentation

Detailed information on endpoints, request/response formats, and error handling are available in the `docs` directory.

## Contributing

We welcome contributions! Fork the repository, raise issues, and submit pull requests. Follow guidelines in `CONTRIBUTING.md`.

## License

MIT license (See `LICENSE` file)

## Support

Open an issue on GitHub for bug reports/feature requests. Join the community discussion on our Discord server (link coming soon).

**Happy building with EKIVAL-Canvas!**
