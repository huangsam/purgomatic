# Purgomatic 💎

A high-performance, Audit-Only tool for tracking and planning large photo/video migrations.

## Philosophy: Home-First Intelligence
Purgomatic is built on the **"Home-First"** principle. Instead of just finding duplicates, it identifies **"Golden Winners"** (files already safely archived in your target Synology/Archive folders) and helps you eliminate the "Toil" (redundant copies on Phones/Thumb drives).

It provides a **Strategic Dashboard** telling you what needs to be moved and what’s already safe.

### Key Features
- **Zero CGO**: Pure Go using `modernc.org/sqlite`.
- **Consolidated Audit**: Single command for scanning and reporting.
- **Hardware Scaled**: Dynamic concurrency via `runtime.NumCPU()`.
- **Multi-Point Hashing**: Sparse SHA-256 sampling for speed and safety.
- **Stat-First Sync**: Skips unchanged files using metadata.
- **Worst Offenders**: Flags top 3 largest files per year.

## Getting Started

### 1. Installation
Purgomatic is a standard Go tool. Once you've pushed this to GitHub:
```bash
go install github.com/[your-user]/purgomatic@latest
```

### 2. Define Scan Targets
Create `scans.json` to tell Purgomatic where your folders are:
```json
[
  { "source": "Synology", "path": "/Volumes/Synology/Photos" },
  { "source": "Local", "path": "/Users/sam/Pictures/Imports" },
  { "source": "Phones", "path": "/Volumes/iPhone_Backup" }
]
```

### 3. The Audit Lifecycle
Using the `Makefile`:
```bash
make init   # Initialize the Multi-Host SQLite database
make audit  # Perform scan and generate global Insight dashboard
```
