# LFX V2 Committee Service

This repository contains the source code for the LFX v2 platform committee service.

## Overview

The LFX v2 Committee Service is a RESTful API service that manages committees within the Linux Foundation's LFX platform. It provides endpoints for creating, reading, updating, and deleting committees with built-in authorization and audit capabilities. Committees are associated with projects and can have hierarchical structures with parent-child relationships.

## File Structure

```bash
├── .github/                        # Github files
│   └── workflows/                  # Github Action workflow files
├── charts/                         # Helm charts for running the service in kubernetes
├── cmd/                            # Services (main packages)
│   └── committee-api/              # Committee service code
│       ├── design/                 # API design specifications (Goa)
│       ├── service/                # Service implementation
│       ├── main.go                 # Application entry point
│       └── http.go                 # HTTP server setup
├── gen/                            # Generated code from Goa design
├── internal/                       # Internal service packages
│   ├── domain/                     # Domain logic layer (business logic)
│   │   ├── model/                  # Domain models and entities
│   │   └── port/                   # Repository and service interfaces
│   ├── service/                    # Service logic layer (use cases)
│   ├── infrastructure/             # Infrastructure layer
│   │   ├── auth/                   # Authentication implementations
│   │   ├── nats/                   # NATS storage implementation
│   │   └── mock/                   # Mock implementations for testing
│   └── middleware/                 # HTTP middleware components
└── pkg/                            # Shared packages
```

## Key Features

- **RESTful API**: Full CRUD operations for committee management
- **Committee Hierarchies**: Support for parent-child committee relationships
- **Project Integration**: Committees are associated with projects for organizational structure
- **Clean Architecture**: Follows clean architecture principles with clear separation of domain, service, and infrastructure layers
- **NATS Storage**: Uses NATS key-value buckets for persistent committee data storage
- **Authorization**: JWT-based authentication with Heimdall middleware integration
- **OpenFGA Support**: Fine-grained authorization control for committee access (configurable)
- **Health Checks**: Built-in `/livez` and `/readyz` endpoints
- **Request Tracking**: Automatic request ID generation and propagation
- **Structured Logging**: JSON-formatted logs with contextual information
- **Committee Settings**: Configurable voting, membership, and access control settings

## Development

To contribute to this repository:

1. Fork the repository
2. Commit your changes to a feature branch in your fork. Ensure your commits
   are signed with the [Developer Certificate of Origin
   (DCO)](https://developercertificate.org/).
   You can use the `git commit -s` command to sign your commits.
3. Ensure the chart version in `charts/lfx-v2-committee-service/Chart.yaml` has been
   updated following semantic version conventions if you are making changes to the chart.
4. Submit your pull request

For detailed development instructions, including local setup, testing, and API development guidelines, see the [Committee API README](cmd/committee-api/README.md).

## License

Copyright The Linux Foundation and each contributor to LFX.

This project’s source code is licensed under the MIT License. A copy of the
license is available in `LICENSE`.

This project’s documentation is licensed under the Creative Commons Attribution
4.0 International License \(CC-BY-4.0\). A copy of the license is available in
`LICENSE-docs`.