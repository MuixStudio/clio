# What is Clio

A modern microservices platform built with Go and Next.js, featuring authentication, user management, and system administration capabilities.

## Quick Start

### Running with Docker

```bash
# Build all services
make build-all

# Build and push to registry
make build-push-all
```

### Running Locally

```bash
# Install dependencies
go mod download

# Run user service
go run services/user/main.go

# Run web service
go run services/web/cmd/run/main.go
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

---

Built with ❤️ by [MuixStudio](https://github.com/muixstudio)