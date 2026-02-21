# 🔍 Infrastructure Audit Report

**Date:** 2025-02-21  
**Auditor:** Claude (AI Assistant)  
**Scope:** `infra/` directory  
**Constraints:** $3/month budget, personal project, security-conscious  
**Context:** [MASTERPLAN.md](./MASTERPLAN.md), [deploy-serverless.yml](../.github/workflows/deploy-serverless.yml)

---

## 📊 Executive Summary

The current infrastructure is **functional but has critical gaps** in both security and cost control. With targeted refactoring, we can achieve **$0-1/month** actual costs while significantly improving security posture.

| Category | Status | Priority |
|----------|--------|----------|
| Security | ⚠️ Medium Risk | High |
| FinOps | ⚠️ Optimizable | High |
| Architecture | ✅ Good Foundation | Medium |
| Compliance | ✅ Acceptable | Low |

---

## 🔒 Security Findings

### 1. IAM Bootstrap Policy Too Broad
**Location:** `infra/core/iam/main.tf:39-50`

**Issue:** The Terraform deployment policy grants wildcard access to critical services:
```hcl
Action = [
  "s3:*",
  "cloudfront:*",
  "lambda:*",
  "iam:*",        # ⚠️ Can create admin users
  "ssm:*",
  # ...
]
Resource = "*"    # ⚠️ Applies to ALL resources
```

**Risk:** If GitHub Actions OIDC token is compromised, attacker gains full account access.

**Remediation:**
```hcl
# Scope to specific resources
Resource = [
  "arn:aws:s3:::trackstack-*",
  "arn:aws:lambda:${var.aws_region}:${local.account_id}:function:trackstack-*",
  "arn:aws:iam::${local.account_id}:role/trackstack-*",
  # ... etc
]
```

### 2. Missing WAF Protection
**Location:** `infra/environments/serverless/cloudfront.tf:72`

**Issue:** CloudFront distribution has no Web Application Firewall.

**Risk:** Vulnerable to:
- DDoS attacks (can spike Lambda costs)
- SQL injection
- XSS attacks
- Bot scraping

**Remediation:** Add AWS WAFv2 with AWS Managed Rules (free tier covers 10M requests/month).

### 3. Lambda Environment Variable Secrets
**Location:** `infra/environments/serverless/lambda.tf:38-54`

**Issue:** Sensitive values stored in Lambda environment variables:
- Origin verification secret
- Turso database tokens (as SSM paths)

**Risk:** Anyone with `lambda:GetFunction` can read these.

**Remediation:** Fetch secrets from SSM/Secrets Manager at runtime using Lambda extensions.

### 4. No Resource Policies on Lambda Function URL
**Location:** `infra/environments/serverless/lambda.tf:66-78`

**Issue:** Lambda Function URL uses IAM auth but lacks resource-based policy restrictions.

**Risk:** Potential for confused deputy attacks.

**Remediation:** Add explicit `aws_lambda_permission` with `source_arn` conditions.

---

## 💰 FinOps Findings

### Current Cost Estimate (Serverless Environment)

| Service | Estimated Monthly Cost |
|---------|----------------------|
| CloudFront | ~$0.50-1.00 (1GB transfer) |
| Lambda | ~$0.00-0.50 (1M requests) |
| S3 | ~$0.10-0.20 (storage + requests) |
| CloudWatch Logs | ~$0.20-0.50 (retention) |
| **Total** | **~$1.00-2.50/month** |

**Within $3 budget** ✅ but can be optimized further.

### 1. CloudFront Price Class
**Location:** `infra/environments/serverless/cloudfront.tf:72`

**Issue:** No explicit `price_class` = defaults to all edge locations (most expensive).

**Impact:** ~30% unnecessary cost for global edge locations not needed for a personal project.

**Remediation:**
```hcl
resource "aws_cloudfront_distribution" "ssr" {
  price_class = "PriceClass_100"  # NA/EU only
  # ...
}
```

**Savings:** ~$0.30-0.50/month

### 2. Lambda Memory Over-provisioned
**Location:** `infra/environments/serverless/lambda.tf:35`

**Issue:** `memory_size = 1024` MB for Astro SSR.

**Impact:** Likely using only 300-500MB, paying for unused capacity.

**Remediation:** 
1. Start with 512MB
2. Monitor with CloudWatch metrics
3. Use AWS Lambda Power Tuning tool to find optimal setting

**Savings:** ~30-50% on Lambda costs (~$0.10-0.25/month)

### 3. S3 Versioning on Assets Bucket
**Location:** `infra/environments/serverless/s3.tf:11-16`

**Issue:** Assets bucket has versioning enabled.

**Impact:** Immutable hashed assets (e.g., `/_astro/index.abc123.js`) are already unique per build. Versioning creates redundant storage.

**Remediation:** Disable versioning on assets bucket:
```hcl
resource "aws_s3_bucket_versioning" "assets" {
  bucket = aws_s3_bucket.assets.id
  versioning_configuration {
    status = "Disabled"  # Assets are immutable
  }
}
```

**Keep versioning on artifacts bucket** (deployment rollbacks).

### 4. CloudWatch Log Retention
**Location:** `infra/environments/serverless/lambda.tf:19-22`

**Issue:** 14-day retention for Lambda logs.

**Impact:** Unnecessary storage for a personal project.

**Remediation:** Reduce to 3 days:
```hcl
resource "aws_cloudwatch_log_group" "ssr" {
  retention_in_days = 3
}
```

**Savings:** ~$0.10-0.20/month

### 5. Missing Cost Controls

**Critical Gap:** No billing alarms or spend limits.

**Remediation:** Add AWS Budgets + CloudWatch billing alarm:
```hcl
resource "aws_cloudwatch_metric_alarm" "billing" {
  alarm_name          = "trackstack-monthly-budget"
  comparison_operator = "GreaterThanThreshold"
  evaluation_periods  = "1"
  metric_name         = "EstimatedCharges"
  namespace           = "AWS/Billing"
  period              = "86400"  # Daily check
  statistic           = "Maximum"
  threshold           = "3"      # $3 USD
  alarm_description   = "Alert when monthly charges exceed $3"
  alarm_actions       = [aws_sns_topic.billing_alerts.arn]
}
```

### 6. Lambda Concurrency Limit
**Location:** `infra/environments/serverless/lambda.tf:24-64`

**Issue:** No reserved concurrency limit.

**Risk:** DDoS or runaway process could generate unlimited Lambda invocations.

**Remediation:** Set conservative limit:
```hcl
resource "aws_lambda_function" "ssr" {
  reserved_concurrent_executions = 10  # Reasonable for personal use
  # ...
}
```

---

## 🏗️ Architecture Findings

### 1. Inline Implementation vs. Modular Design
**Location:** `infra/environments/serverless/*.tf`

**Issue:** All resources defined inline. [MASTERPLAN.md](./MASTERPLAN.md) mentions reusable modules but none exist yet.

**Impact:** 
- Code duplication when adding ECS/EKS environments
- Harder to maintain consistency
- Violates DRY principle

**Recommended Structure:**
```
infra/
├── modules/
│   ├── lambda-api/           # Lambda + Function URL + IAM role
│   ├── static-hosting/       # S3 + CloudFront + OAC
│   ├── security-headers/     # CloudFront response headers policy
│   ├── waf-core/             # WAF with basic rules
│   └── cost-guardrails/      # Budgets + billing alarms
├── environments/
│   └── serverless/
│       ├── main.tf           # Just wire modules together
│       ├── variables.tf
│       └── terraform.tfvars
```

### 2. Mixed Concerns in Single Directory
**Location:** `infra/environments/serverless/`

**Issue:** Contains both infrastructure (Terraform) and deployment configuration (via SSM parameters).

**Recommendation:** Keep deployment config (bucket names, distribution IDs) in a separate config file or use Terraform outputs more explicitly.

### 3. Hardcoded Runtime Configuration
**Location:** `infra/environments/serverless/lambda.tf:38-54`

**Issue:** Lambda environment variables hardcode SSM parameter paths:
```hcl
TURSO_USERS_URL = "/${local.ssm_prefix}/runtime/TURSO_USERS_URL"
```

**Recommendation:** Pass as map variable:
```hcl
variable "turso_ssm_paths" {
  type = map(string)
  default = {}
}

environment {
  variables = var.turso_ssm_paths
}
```

---

## ✅ What's Working Well

1. **Hexagonal Architecture Alignment:** Lambda and CloudFront setup supports the goal of running same code in different environments.

2. **Origin Verification:** CloudFront custom headers prevent direct Lambda access.

3. **OAC Usage:** Modern Origin Access Control for S3 (better than legacy OAI).

4. **Security Headers:** Comprehensive response headers policy (HSTS, XSS, CSP, etc.).

5. **Least Privilege S3:** Buckets have public access blocked, ownership enforced.

6. **GitHub Actions OIDC:** Proper use of short-lived credentials (no long-term access keys).

---

## 🎯 Prioritized Remediation Roadmap

### Phase 1: Critical (Do This Week)

- [ ] **Add Billing Alarm** - Prevent surprise bills
  ```bash
  terraform apply -target=aws_cloudwatch_metric_alarm.billing
  ```

- [ ] **Scope IAM Policy** - Reduce blast radius
  - Edit `infra/core/iam/main.tf`
  - Replace wildcard resources with specific ARNs

- [ ] **Add WAF** - Basic DDoS/injection protection
  - Use AWS WAFv2 with free managed rules

### Phase 2: High Priority (Next 2 Weeks)

- [ ] **CloudFront Price Class** - Save 30% immediately
  - Add `price_class = "PriceClass_100"`

- [ ] **Lambda Memory Tuning** - Right-size compute
  - Reduce to 512MB initially
  - Run AWS Lambda Power Tuning

- [ ] **Create Reusable Modules** - Enable ECS/EKS labs
  - Extract `lambda-api` module
  - Extract `static-hosting` module

### Phase 3: Optimization (Ongoing)

- [ ] **Disable S3 Versioning on Assets** - Remove waste
- [ ] **Reduce Log Retention** - 14 days → 3 days
- [ ] **Lambda Concurrency Limit** - Prevent runaway costs
- [ ] **Runtime Secret Fetching** - Move secrets out of env vars

---

## 📈 Expected Outcomes

After Phase 1 + 2:
- **Security Score:** Medium → High
- **Monthly Cost:** $1.00-2.50 → $0.50-1.50
- **Maintainability:** Improved via modularization
- **Risk:** Significantly reduced blast radius

---

## 🔗 Related Documents

- [MASTERPLAN.md](./MASTERPLAN.md) - Architecture strategy
- [deploy-serverless.yml](../.github/workflows/deploy-serverless.yml) - CI/CD pipeline
- AWS Free Tier Limits: https://aws.amazon.com/free/
- AWS WAF Pricing: https://aws.amazon.com/waf/pricing/

---

*Audit completed on 2025-02-21. Review quarterly or when adding new environments (ECS/EKS).*
