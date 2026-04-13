provider "aws" {
  region = "us-east-2"
}

# Define the AWS S3 bucket
resource "aws_s3_bucket" "log_storage" {
  bucket = "juanes-cloud-logs-${random_id.suffix.hex}"
}

resource "random_id" "suffix" {
  byte_length = 4
}

# Output: created bucket name
output "bucket_name" {
  value = aws_s3_bucket.log_storage.id
}

resource "aws_ecr_repository" "analyzer_repo" {
  name                 = "cloud-log-analyzer"
  image_tag_mutability = "MUTABLE"

  image_scanning_configuration {
    # Security: scan images searching vulnerabilities
    scan_on_push = true
  }
}

output "ecr_url" {
  value = aws_ecr_repository.analyzer_repo.repository_url
}
