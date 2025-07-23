

## Configuration

- Create a `.env` file in the root directory based on the `.env.example` file to run the server locally.
- In the `rarible-helm/templates/` directory, create a `secret.yaml` file based on `.secret.example.yaml` for deploying to a Kubernetes cluster (local or remote).

## Running the Application

You can start the application using Docker via Makefile (from Docker hub):
```bash
make start
```

or with Docker Compose(from Docker hub):
```bash
docker-compose up --build -d
```

Alternatively, you can build the Docker image locally with:
```bash
docker build -t <name>:<tag> .
```

Running the Docker Image with .env on Port 8080
```bash
docker run --env-file .env -p 8080:8080 <name>:<tag>
```

## Install the Helm chart to your Kubernetes cluster
```bash
helm install <release-name> rarible-helm
```

## Running Tests

To run automated tests for the Client and Service, use:

```bash
make test
```

Automated tests are also run automatically on every push via GitHub Actions workflows.

## API Documentation

After starting the application, the Swagger UI documentation will be available at:
```bash
http://localhost:8080/swagger/index.html#/
```

The documentation includes:
- Detailed endpoint descriptions
- Request/response data models
- Model schemas
