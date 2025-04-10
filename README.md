# Countries App

This is a simple microservice application built with Go that consists of two services:

1. **Countries Service**: Reads country data from a CSV file and exposes it via a REST API
2. **Client Service**: Fetches country data from the Countries Service and displays it in an HTML table

## Project Structure

```
countries-app/
├── countries-service/      # First service - Country data API
│   ├── main.go
│   ├── Dockerfile
│   └── countries.csv       # Source data
├── client-service/         # Second service - Web UI
│   ├── main.go
│   └── Dockerfile
├── docker-compose.yml      # Local development setup
└── helm-umbrella-chart/    # Kubernetes deployment with Helm
    ├── Chart.yaml
    ├── values.yaml
    └── charts/
        ├── countries-service/
        │   ├── Chart.yaml
        │   ├── values.yaml
        │   └── templates/
        │       ├── configmap.yaml
        │       ├── deployment.yaml
        │       └── service.yaml
        └── client-service/
            ├── Chart.yaml
            ├── values.yaml
            └── templates/
                ├── deployment.yaml
                ├── service.yaml
                └── ingress.yaml
```

## Running Locally with Docker Compose

1. Build and start the services:
   ```
   docker-compose up --build
   ```

2. Access the client service in your browser at:
   ```
   http://localhost:8081
   ```

3. You can also directly access the countries API at:
   ```
   http://localhost:8080/countries
   ```

## Deploying with Kubernetes and Helm

### Prerequisites
- Kubernetes cluster
- Helm installed
- Docker registry (to push images)

### Build and Push Docker Images

```bash
# Build countries-service image
docker build -t your-registry/countries-service:latest ./countries-service
docker push your-registry/countries-service:latest

# Build client-service image
docker build -t your-registry/client-service:latest ./client-service
docker push your-registry/client-service:latest
```

### Deploy with Helm

1. Update the image repositories in `helm-umbrella-chart/values.yaml` to match your registry:

```yaml
countries-service:
  image:
    repository: your-registry/countries-service
    
client-service:
  image:
    repository: your-registry/client-service
```

2. Install the chart:

```bash
helm install countries-app ./helm-umbrella-chart
```

3. To access the application:

```bash
# If using a NodePort service
kubectl port-forward svc/countries-app-client-service 8081:8081

# If ingress is enabled (update your /etc/hosts file to point countries.local to your cluster IP)
http://countries.local
```

## Development

To modify the country data, edit the `countries-service/countries.csv` file.