# Committee API

This directory contains the Committee API service. The service provides comprehensive committee management functionality including:

- It serves HTTP requests via Traefik to perform CRUD operations on committee data
- Manages committee base information (name, category, description, settings, etc.)
- Handles committee settings including voting configurations, member management, and access controls
- Integrates with project services to ensure committee-project relationships

Applications with a BFF should use the REST API with HTTP requests to perform the needed operations on committees, while other resource API services can communicate with this service as needed.

This service contains the following API endpoints:

- `/committees`
  - `POST`: create a new committee with base information and settings

## File Structure

```bash
├── design/                         # Goa design files
│   ├── committee.go                # Goa committee service specification
│   └── type.go                     # Goa data types and models
├── service/                        # Service implementation (presentation layer)
│   ├── committee_service.go        # Committee service implementation
│   ├── error.go                    # Error handling utilities
│   └── providers.go                # Dependency injection providers
├── main.go                         # Application startup and dependency injection
├── http.go                         # HTTP server setup and configuration
└── README.md                       # This documentation

# Dependencies from internal/ packages:
# - internal/service/              # Business logic and use case orchestration
# - internal/domain/               # Domain models, ports, and business rules
# - internal/infrastructure/       # Infrastructure implementations (NATS storage, Auth, Messaging)
# - internal/middleware/           # HTTP middleware components
```

## Architecture

This service follows clean architecture principles with clear separation of concerns:

### Layers

1. **Presentation Layer** (`cmd/committee-api/`)
   - `committeeServicesrvc` struct implements the Goa-generated service interface
   - HTTP endpoint handlers for committee operations (`service/committee_service.go`)
   - HTTP server setup and configuration (`http.go`)
   - Dependency injection and startup (`main.go`)

2. **Service/Use Case Layer** (`internal/service/`)
   - `CommitteeWriter` orchestrates committee creation and management
   - Contains business logic for committee operations
   - Validates business rules and coordinates between domain and infrastructure

3. **Domain Layer** (`internal/domain/`)
   - Domain models (`model/`)
   - Port interfaces (`port/`)
   - Business rules and domain-specific validation

4. **Infrastructure Layer** (`internal/infrastructure/`)
   - NATS storage implementation (`nats/`)
   - JWT authentication implementation (`auth/`)
   - Mock implementations for testing (`mock/`)
   - Messaging infrastructure

### Key Benefits

- **Storage Independence**: Can switch from NATS to PostgreSQL without changing business logic
- **Testability**: Each layer can be tested in isolation using comprehensive mocks
- **Maintainability**: Clear separation of concerns and dependency direction
- **Scalability**: Support for committee hierarchies and complex organizational structures
- **Integration**: Seamless integration with project services and external authentication systems

## Development

### Prerequisites

- [**Go**](https://go.dev/): the service is built with the Go programming language [[Install](https://go.dev/doc/install)]
- [**Kubernetes**](https://kubernetes.io/): used for deployment of resources [[Install](https://kubernetes.io/releases/download/)]
- [**Helm**](https://helm.sh/): used to manage kubernetes applications [[Install](https://helm.sh/docs/intro/install/)]
- [**NATS**](https://docs.nats.io/): used to communicate with other LFX V2 services [[Install](https://docs.nats.io/running-a-nats-service/introduction/installation)]
- [**GOA Framework**](https://goa.design/): used for API code generation

#### GOA Framework

Follow the [GOA installation guide](https://goa.design/docs/2-getting-started/1-installation/) to install GOA:

```bash
go install goa.design/goa/v3/cmd/goa@latest
```

Verify the installation:

```bash
goa version
```

### Building and Development

#### 1. Generate Code

The service uses GOA to generate API code from the design specification. Run the following command to generate all necessary code:

```bash
make apigen

# or directly run the "goa gen" command
goa gen github.com/linuxfoundation/lfx-v2-committee-service/cmd/committee-api/design
```

This command generates:

- HTTP server and client code
- OpenAPI specification
- Service interfaces and types
- Transport layer implementations

#### 2. Set up resources and external services

The service relies on some resources and external services being spun up prior to running this service.

- [NATS service](https://docs.nats.io/): ensure you have a NATS server instance running and set the `NATS_URL` environment variable with the URL of the server

    ```bash
    export NATS_URL=nats://lfx-platform-nats.lfx.svc.cluster.local:4222
    ```

- [NATS key-value bucket](https://docs.nats.io/nats-concepts/jetstream/key-value-store): once you have a NATS service running, you need to create buckets used by the committee service.

    ```bash
    # if using the nats cli tool
    nats kv add committees --history=20 --storage=file --max-value-size=10485760 --max-bucket-size=1073741824
    nats kv add committee-settings --history=20 --storage=file --max-value-size=10485760 --max-bucket-size=1073741824
    ```

#### 3. Export environment variables

|Environment Variable Name|Description|Default|Required|
|-----------------------|--------------------|-----------|-----|
|PORT|the port for http requests to the committee service API|8080|false|
|NATS_URL|the URL of the nats server instance|nats://localhost:4222|false|
|LOG_LEVEL|the log level for outputted logs|info|false|
|LOG_ADD_SOURCE|whether to add the source field to outputted logs|false|false|
|JWKS_URL|the URL to the endpoint for verifying ID tokens and JWT access tokens||false|
|AUDIENCE|the audience of the app that the JWT token should have set - for verification of the JWT token|lfx-v2-committee-service|false|
|JWT_AUTH_DISABLED_MOCK_LOCAL_PRINCIPAL|a mocked auth principal value for local development (to avoid needing a valid JWT token)||false|

#### 4. Development Workflow

1. **Make design or implementation changes**: Edit files in the `design/` directory for design changes, and edit the other files for implementation changes.

2. **Regenerate code**: Run `make apigen` after design changes

3. **Build the service**:

   ```bash
   make build
   ```

4. **Run the service**:

   ```bash
   make run

   # or run with debug logs enabled
   make debug

   # or run with the go command to set custom flags
   # -bind string   interface to bind on (default "*")
   # -d          enable debug logging (default false)
   # -p    string   listen port (default "8080")
   go run
   ```

   Once the service is running, make a request to the `/livez` endpoint to ensure that the service is alive.

   ```bash
    curl http://localhost:8080/livez
   ```

   You should get a 200 status code response with a text/plain content payload of `OK`.

5. **Run tests**:

   ```bash
   make test

   # or run go test to set custom flags
   go test . -v
   ```

6. **Lint the code**

   ```bash
   # From the root of the directory, run megalinter (https://megalinter.io/latest/mega-linter-runner/) to ensure the code passes the linter checks. The CI/CD has a check that uses megalinter.
   npx mega-linter-runner .
   ```

7. **Docker build + K8**

    ```bash
    # Build the dockerfile (from the root of the repo)
    docker build -t lfx-v2-committee-service:<release_number> .

    # Install the helm chart for the service into the lfx namespace (from the root of the repo)
    helm install lfx-v2-committee-service ./charts/lfx-v2-committee-service/ -n lfx

    # Once you have already installed the helm chart and need to just update it, use the following command (from the root of the repo):
    helm upgrade lfx-v2-committee-service ./charts/lfx-v2-committee-service/ -n lfx

    # Check that the REST API is accessible by hitting the `/livez` endpoint (you should get a response of OK if it is working):
    #
    # Note: replace the hostname with the host from ./charts/lfx-v2-committee-service/ingressroute.yaml
    curl http://lfx-api.k8s.orb.local/livez
    ```
### Authorization with OpenFGA

When deployed via Kubernetes, the committee service uses OpenFGA for fine-grained authorization control. The authorization is handled by Heimdall middleware before requests reach the service.

#### Configuration

OpenFGA authorization is controlled by the `openfga.enabled` value in the Helm chart:

```yaml
# In values.yaml or via --set flag
openfga:
  enabled: true  # Enable OpenFGA authorization (default)
  # enabled: false  # Disable for local development only
```

#### Local Development

For local development without OpenFGA:

1. Set `openfga.enabled: false` in your Helm values
2. All requests will be allowed through (after JWT authentication)
3. **Warning**: Never disable OpenFGA in production environments

### Add new API endpoints

Note: follow the [Development Workflow](#4-development-workflow) section on how to run the service code

1. **Update design files**: Edit the committee design file in `design/committee.go` to include specification of the new endpoint with all of its supported parameters, responses, and errors, etc.
2. **Regenerate code**: Run `make apigen` after design changes to generate the new Goa interfaces and types
3. **Implement code**: Implement the new endpoint in `service/` following the existing patterns. Add the necessary business logic to the use case layer in `internal/service/` if needed. Include comprehensive tests for the new endpoint.
4. **Update heimdall ruleset**: Ensure that `/charts/lfx-v2-committee-service/templates/ruleset.yaml` has the route and method for the endpoint set so that authentication is configured when deployed. If the endpoint modifies data (PUT, DELETE, PATCH), consider adding OpenFGA authorization checks in the ruleset for proper access control

### Adding New Business Logic

For complex committee operations that require multiple steps or external service integration:

1. **Extend Use Cases**: Add new methods to the `CommitteeWriter` interface in `internal/service/`
2. **Update Domain Models**: Modify committee models in `internal/domain/model/` if new data structures are needed  
3. **Extend Port Interfaces**: Update port interfaces in `internal/domain/port/` to support new storage or external service operations
4. **Implement Infrastructure**: Add concrete implementations in `internal/infrastructure/` for new external service integrations
5. **Add Tests**: Create comprehensive unit tests with mocks for all new components