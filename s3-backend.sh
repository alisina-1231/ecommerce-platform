#! bin/bash

aws s3api create-bucket \
  --bucket terraform-test-bucket-1231 \
  --region us-east-1

  aws s3api put-bucket-versioning \
  --bucket terraform-test-bucket-1231 \
  --versioning-configuration Status=Enabled

  aws s3api put-bucket-encryption \
  --bucket terraform-test-bucket-1231 \
  --server-side-encryption-configuration '{
    "Rules": [
      {
        "ApplyServerSideEncryptionByDefault": {
          "SSEAlgorithm": "AES256"
        }
      }
    ]
  }'

  aws s3api put-public-access-block \
  --bucket terraform-test-bucket-1231 \
  --public-access-block-configuration \
  BlockPublicAcls=true,IgnorePublicAcls=true,BlockPublicPolicy=true,RestrictPublicBuckets=true