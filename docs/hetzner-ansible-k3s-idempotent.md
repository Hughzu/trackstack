# Hetzner + Ansible + K3s (Idempotent)

This is a practical, idempotent provisioning flow for a single VPS on Hetzner using Ansible. It keeps the first deployment fast while leaving room for a “golden path” later.

## Goals
- Provision a Hetzner VPS with Ansible (idempotent)
- Harden SSH and configure firewall
- Install K3s
- Keep Traefik ingress (fastest)
- Install cert-manager for Let’s Encrypt TLS
- Route traffic via Ingress on 80/443

## Overview
1. Create server via Hetzner API (Ansible `hcloud` module).
2. Bootstrap server (users, SSH hardening, UFW).
3. Install K3s.
4. Install cert-manager.
5. Configure DNS A record for your domain.
6. Apply app Ingress with TLS.

## Prerequisites
- A domain name you control
- Hetzner API token
- SSH key for server access
- Ansible installed locally
- Python packages: `hcloud` (Ansible module dependency)

## Suggested Repo Layout
```
infra/
  ansible/
    inventories/
      prod/
        hosts.ini
        group_vars/
          all.yml
    playbooks/
      provision.yml
      bootstrap.yml
      k3s.yml
      cert-manager.yml
    roles/
      common/
      k3s/
      cert_manager/
```

## Inventory Example
`infra/ansible/inventories/prod/hosts.ini`
```
[k3s]
your-server ansible_host=<PUBLIC_IP> ansible_user=root
```

`infra/ansible/inventories/prod/group_vars/all.yml`
```
domain_name: app.example.com
letsencrypt_email: you@example.com
k3s_version: v1.29.4+k3s1
ssh_user: deploy
ssh_public_key: "ssh-ed25519 AAAA..."
```

## Playbooks

### 1) provision.yml (Create VPS)
Uses Hetzner API to create the server if it doesn’t exist.

Key tasks:
- Create server with `hcloud_server`
- Attach SSH key
- Output server IP for inventory update

Idempotency:
- `hcloud_server` is declarative: re-run keeps state

### 2) bootstrap.yml (Base OS)
Key tasks:
- Update packages
- Create `deploy` user
- Set authorized keys
- Harden SSH (disable password, allow key only)
- Configure UFW (22, 80, 443)

Idempotency:
- Ansible modules are declarative; re-run safe

### 3) k3s.yml (K3s install)
Key tasks:
- Install K3s server (pinned version)
- Ensure systemd service enabled

Idempotency:
- Check if `k3s` exists and skip install if present
- Verify node Ready

### 4) cert-manager.yml
Key tasks:
- Install cert-manager via Helm or manifest
- Apply `ClusterIssuer` for Let’s Encrypt HTTP-01

Idempotency:
- Helm/manifests applied repeatedly are safe

## DNS
Create an A record:
- `app.example.com` -> `<PUBLIC_IP>`

## Ingress (App)
Your app should be exposed via an Ingress resource with TLS:
- `host: app.example.com`
- `tls.secretName: app-tls`
- `cert-manager.io/cluster-issuer: letsencrypt-prod`

## Minimal Execution Flow
From your machine:
```
ansible-playbook -i infra/ansible/inventories/prod/hosts.ini infra/ansible/playbooks/provision.yml
ansible-playbook -i infra/ansible/inventories/prod/hosts.ini infra/ansible/playbooks/bootstrap.yml
ansible-playbook -i infra/ansible/inventories/prod/hosts.ini infra/ansible/playbooks/k3s.yml
ansible-playbook -i infra/ansible/inventories/prod/hosts.ini infra/ansible/playbooks/cert-manager.yml
```

## CI/CD Later
Once stable locally, move to CI with:
- Hetzner API token in CI secrets
- SSH private key in CI secrets
- Idempotent playbooks on push or manual trigger

## Notes
- Keep ports 80/443 only; do not expose arbitrary service ports.
- Keep Traefik for speed; swap to NGINX later if needed.
- Use DNS + cert-manager for TLS; avoid IP-only HTTPS.
