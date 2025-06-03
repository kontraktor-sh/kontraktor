# Kontraktor – Architecture Design Document

> **Version:** 0.2 (June 3 2025)

---

## 1  Purpose & Vision
Kontraktor is a **builder‑helper** that unifies environment configuration, secret aggregation and developer task automation. Inspired by [taskfile.dev](https://taskfile.dev) but extended with a *centralised configuration library* and *remote‑vault integration*, it offers:

1. **Central Configuration Library** – version‑controlled, queryable configuration across many projects and git repositories, stored in DuckDB (non‑secret only).
2. **Secret Aggregation Layer** – read‑only connectors to Azure Key Vault, HashiCorp Vault and AWS Secrets Manager (no secret persisted in Kontraktor).
3. **Task Runner** – local execution of Bash (macOS/Linux) and, later, PowerShell Core routines through an expressive YAML syntax (`taskfile.ktr.yml`).





