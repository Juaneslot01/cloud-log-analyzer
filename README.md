# 🚀 Cloud-Native Log Ingestor & Infrastructure

**High-performance Backend architecture built with Go and gRPC, automated via DevOps best practices.**

This project demonstrates the implementation of a high-availability data ingestion system designed to handle massive information streams (from IoT/DDS environments) and persist them into the AWS Cloud using a decoupled and scalable architecture.

## 🏗️ Backend Engineering Highlights

* **Concurrent Processing (Worker Pool Pattern):** Implementation of a producer/consumer model in **Go** using `channels` and `goroutines`. This allows the gRPC server to accept thousands of requests per second without blocking the main thread, significantly optimizing throughput.
* **Efficient Communication with gRPC:** Leverages **Protocol Buffers** to define a strictly typed and efficient API contract, reducing payload size and improving performance compared to traditional REST architectures.
* **Environment-Agnostic Design:** Follows **Twelve-Factor App** principles. The business logic is entirely decoupled from configuration through dynamic injection of environment variables and runtime secrets.

## 🛠️ DevOps & Infrastructure as Code (IaC)

* **Kubernetes Orchestration:** Managed deployment on a **Kubernetes (k3d)** cluster, utilizing `Deployments` for horizontal scaling and `Services` (LoadBalancer) for traffic exposure.
* **Infrastructure Automation (Terraform):** The entire cloud infrastructure (AWS S3, ECR, and IAM) is codified using **Terraform**, enabling the creation of identical and reproducible dev/prod environments in minutes.
* **CI/CD Pipeline:** Automated continuous delivery flow via **GitHub Actions**. The pipeline handles Docker image builds, pushes to Amazon ECR, and manages versioning through commit hashes to ensure full traceability.
* **Containerization:** Uses Docker to package the application, ensuring a consistent runtime environment from local Linux development to AWS production.

## 📊 Technical Stack

| Category | Technologies |
| :--- | :--- |
| **Backend** | Go (Golang), gRPC, Protobuf |
| **Infrastructure** | Terraform, AWS (S3, ECR), IAM |
| **Orchestration** | Kubernetes (k8s), Docker, k3d |
| **Automation** | GitHub Actions, Git |
| **Environment** | Arch-based Linux, Neovim/Ghostty |

## ⚙️ Development Workflow

1.  **Local Dev:** Building low-latency microservices in Go.
2.  **IaC:** Provisioning cloud storage and permissions via `terraform apply`.
3.  **Pipeline:** Code is pushed to GitHub; Actions build and store the container image.
4.  **Deploy:** Kubernetes pulls the latest image and manages the Pod lifecycle.
5.  **Observability:** Monitoring logs and container states to ensure system stability.


## 🛠️ Getting Started

### Prerequisites

Ensure you have the following tools installed:
* **Go** (1.21+)
* **Docker** & **k3d**
* **Terraform**
* **AWS CLI** (configured with valid credentials)
* **grpcurl** (for testing)

### 1. Infrastructure Provisioning (AWS)

First, spin up the required cloud resources (S3, ECR, IAM):
```bash
cd terraform
terraform init
terraform apply
```

### 2. Local Cluster Setup

Create a local Kubernetes cluster using k3d:
```bash
k3d cluster create my-cluster -p "50051:50051@loadbalancer"
```

### 3. Deploy to Kubernetes

Apply the secrets and manifests. Make sure to update the image tag in deployment.yaml with your ECR URI:
```bash
# Create secrets for AWS credentials and App API Key
kubectl create secret generic aws-creds \
  --from-literal="access_key=$AWS_ACCESS_KEY_ID" \
  --from-literal="secret_key=$AWS_SECRET_ACCESS_KEY" \
  --from-literal="api_key=your_chosen_api_key"

# Deploy the analyzer
kubectl apply -f k8s/
```
### 4. Testing the Service

Once the pods are running, send a sample log via gRPC:
```bash
grpcurl -plaintext \
  -rpc-header "api-key: your_chosen_api_key" \
  -d '{"service_name": "manual-test", "level": "INFO", "message": "Hello from local terminal"}' \
  localhost:50051 ingestor.LogService/SendLog
```

### Makefile quick commands

**1. To see all commands:**
```bash
make help
```

**2. To initialize cloud and cluster:**
```bash
make infra-up
make cluster-up
```

**3. For configuration of credentials:**
```bash
make secret
```

**4. For Deploying and logs**
```bash
make deploy
make logs
```

**5. For testing the project**
```bash
make test
```
