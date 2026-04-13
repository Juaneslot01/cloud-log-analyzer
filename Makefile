# Variables
CLUSTER_NAME=k3d-cluster
IMAGE_NAME=785894519798.dkr.ecr.us-east-2.amazonaws.com/cloud-log-analyzer:latest
API_KEY=analizador-secreto-2026

.PHONY: help infra-up infra-down cluster-up cluster-down deploy secret logs test

help: ## Show this help message
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

## --- Infrastructure (Terraform) ---
infra-up: ## Initialize and apply terraform changes
	cd terraform && terraform init && terraform apply -auto-approve

infra-down: ## Destroy cloud infrastructure
	cd terraform && terraform destroy -auto-approve

## --- Local Cluster (k3d) ---
cluster-up: ## Create k3d cluster with gRPC port mapping
	k3d cluster create $(CLUSTER_NAME) -p "50051:50051@loadbalancer"

cluster-down: ## Delete local k3d cluster
	k3d cluster delete $(CLUSTER_NAME)

## --- Kubernetes Deployment ---
secret: ## Create Kubernetes secrets for AWS and App
	@read -p "Enter AWS_ACCESS_KEY_ID: " acc; \
	read -p "Enter AWS_SECRET_ACCESS_KEY: " sec; \
	kubectl create secret generic aws-creds \
		--from-literal=access_key=$$acc \
		--from-literal=secret_key=$$sec \
		--from-literal=api_key=$(API_KEY)

deploy: ## Apply K8s manifests
	kubectl apply -f k8s/

restart: ## Force a rollout restart of the deployment
	kubectl rollout restart deployment cloud-log-analyzer

logs: ## Tail application logs
	kubectl logs -l app=cloud-log-analyzer -f --tail=50

## --- Testing ---
test: ## Send a test log using grpcurl
	grpcurl -plaintext \
		-rpc-header "api-key: $(API_KEY)" \
		-d '{"service_name": "makefile-test", "level": "INFO", "message": "Automation is working!"}' \
		localhost:50051 ingestor.LogService/SendLog
