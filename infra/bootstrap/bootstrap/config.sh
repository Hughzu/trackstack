#!/bin/bash

# AWS Configuration
export AWS_REGION="eu-west-1"
export AWS_ACCOUNT_ID="939091506005"
export AWS_PROFILE="${AWS_PROFILE:-trackstack}"

# GitHub Configuration
export GITHUB_REPO="Hughzu/trackstack"
export GITHUB_ORG="Hughzu"
export GITHUB_OIDC_SUBJECT="repo:${GITHUB_REPO}:ref:refs/heads/main"

# Resource Names
export ADMIN_USER_NAME="trackstack-admin"
export TFSTATE_BUCKET="trackstack-tfstate-939091506005"
export TFSTATE_LOCK_TABLE="trackstack-terraform-locks"
export GITHUB_ROLE_NAME="trackstack-github-deploy-role"
export TERRAFORM_POLICY_NAME="trackstack-terraform-deploy-policy"
export OIDC_PROVIDER_URL="token.actions.githubusercontent.com"
