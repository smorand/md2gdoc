.PHONY: build run test clean docker-build docker-push deploy tf-init tf-plan tf-apply tf-destroy

# Variables
PROJECT_ID   ?= your-gcp-project-id
REGION       ?= europe-west1
SERVICE_NAME ?= md2gdoc
IMAGE        ?= $(REGION)-docker.pkg.dev/$(PROJECT_ID)/$(SERVICE_NAME)/$(SERVICE_NAME)
TAG          ?= latest
PORT         ?= 8080

# --- Go ---

build:
	CGO_ENABLED=0 go build -o $(SERVICE_NAME) .

run: build
	PORT=$(PORT) ./$(SERVICE_NAME)

test:
	go test -v ./...

clean:
	rm -f $(SERVICE_NAME)

# --- Docker ---

docker-build:
	docker build -t $(IMAGE):$(TAG) .

docker-push: docker-build
	docker push $(IMAGE):$(TAG)

docker-run:
	docker run --rm -p $(PORT):8080 $(IMAGE):$(TAG)

# --- GCP Auth ---

gcp-auth:
	gcloud auth configure-docker $(REGION)-docker.pkg.dev

# --- Terraform ---

tf-init:
	cd terraform && terraform init

tf-plan:
	cd terraform && terraform plan

tf-apply:
	cd terraform && terraform apply

tf-destroy:
	cd terraform && terraform destroy

# --- Full deploy pipeline ---

deploy: docker-push tf-apply
	@echo "Deployed $(IMAGE):$(TAG)"
	@echo "Endpoint: https://md2gdoc.mcp.scm-platform.org/mcp"
