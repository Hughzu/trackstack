# 🏗️ TrackStack Infrastructure Strategy: The Multi-Tenant Golden Path

**Date:** 2026-02-10
**Context:** Defining the deployment strategy for the "Golden Paths".
**Decision:** Adopt a **Multi-Tenant K3s Cluster** on a single VPS instead of "One VPS per App".

## 1. The Core Strategy: "Shared Metal, Isolated Logic"

To align with the "Product First" goal, we prioritize deployment speed ("Time to Hello World") and resource efficiency over strict physical isolation.

* **Architecture:** A single, robust VPS (Hetzner) running a K3s cluster.
* **Isolation:** Logical isolation via Kubernetes **Namespaces**.
* **Routing:** Dynamic routing via **Ingress Controller** (Traefik) and Wildcard DNS.
* **Outcome:** New apps can be deployed/destroyed in seconds without provisioning new VMs.

## 2. The Tooling Triad

We separate concerns into three distinct layers to ensure idempotency and stability.

| Layer | Tool | Responsibility | Frequency |
| --- | --- | --- | --- |
| **Foundation** | **Terraform** | **"The Hardware"**<br>

<br>- Provisioning the Hetzner VPS.<br>

<br>- Managing Cloudflare DNS records (Wildcard `*.apps.trackstack.dev`).<br>

<br>- Managing Block Storage Volumes (for SQLite persistence). | Once / Rare |
| **Config** | **Ansible** | **"The OS & Runtime"**<br>

<br>- OS Security (UFW, Fail2Ban).<br>

<br>- Installing/Upgrading K3s.<br>

<br>- Injecting `kubeconfig`.<br>

<br>- Automated OS Patching (`unattended-upgrades`). | Monthly / Maintenance |
| **Delivery** | **CLI / Helm** | **"The Application"**<br>

<br>- Interacting with K8s API.<br>

<br>- creating Namespaces & Ingress rules.<br>

<br>- pulling images from GHCR.<br>

<br>- Managing App Lifecycle (Deploy/Kill). | Daily / Continuous |

## 3. Network & Routing Strategy

* **DNS:** A single Wildcard A Record in Cloudflare pointing to the VPS IP.
* `*.apps.trackstack.dev` -> `1.2.3.4`


* **Ingress:** The CLI generates a Kubernetes `Ingress` resource for every deployment.
* Example: `budget-app.apps.trackstack.dev` routes automatically to Service `budget-app` in Namespace `budget`.


* **Zero-Touch Infra:** No Terraform/Ansible runs are required to add a new subdomain.

## 4. Maintenance & OS Lifecycle ("Cattle not Pets")

To mitigate the operational overhead of managing a VPS:

### Routine Maintenance

* **Automated Patching:** Ansible configures `unattended-upgrades` for automatic security patches.
* **Kernel Updates:** Use **Kured** (Kubernetes Reboot Daemon) to watch for `/var/run/reboot-required`, drain nodes, and reboot safely.

### Major Upgrades (The "Reprovision" Pattern)

Instead of `do-release-upgrade`:

1. **Terraform:** Spin up a new VPS with the latest OS image.
2. **Ansible:** Bootstrap K3s on the new node.
3. **Migration:** Detach the **Block Storage Volume** (containing SQLite data) from the old VPS and attach to the new one.
4. **Switch:** Update the Cloudflare IP (Terraform).
5. **Destroy:** Terminate the old VPS.

## 5. Implementation Directives for Agents

When generating code, adhere to these constraints:

* **CLI Idempotency:** The CLI must check for existing K3s contexts before attempting operations. It must handle "adoption" of existing resources.
* **No ClickOps:** All Cloudflare and Hetzner configurations must be defined in `.tf` files.
* **Stateless Compute:** Ensure applications expect SQLite databases to be mounted via Persistent Volume Claims (PVC) linked to the Block Storage, not the ephemeral root disk.
