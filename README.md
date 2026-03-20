# Purgomatic 💎 (v1.1)

**The SRE Photo Librarian** — A high-performance, minimalist Audit-Only tool for tracking and planning large photo/video migrations (91k+ files).

---

## Philosophy: Home-First Intelligence
Purgomatic is built on the **"Home-First"** principle. Instead of just finding duplicates, it identifies **"Golden Winners"** (files already safely archived in your target Synology/Archive folders) and helps you eliminate the "Toil" (redundant copies on Phones/Thumb drives).

It provides a **Strategic Dashboard** telling you what needs to be moved and what’s already safe.

### Key Features
- **100% Pure Go (Zero CGO)**: No external C-dependencies. Built with `modernc.org/sqlite`.
- **Consolidated Audit**: A single `audit` command indexes your library and generates a high-level strategic report.
- **Hardware Optimized**: Scalable concurrency utilizing `runtime.NumCPU() * 2`. Blazing fast on modern multi-core systems (like the Apple M3 Max).
- **Multi-Point Hashing**: Custom SHA-256 sampler (First/Middle/End) for collision-resistant deduplication at scale.
- **Stat-First Sync**: Blazing fast rescans by skipping hashing for unchanged files based on Size/Mtime.
- **Worst Offenders Tracking**: Automatically flags the top 3 largest files per year to highlight migration targets.

---

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
make audit  # Perform scan and generate global SRE Insight dashboard
```

---

## Architectural Specs
| Metric | Specification |
| :--- | :--- |
| **Language** | Go 1.26+ (Modern Concurrency) |
| **Database** | SQLite (CGO-free / Multi-Host) |
| **Concurrency** | Dynamic (`runtime.NumCPU() * 2`) |
| **Hashing** | Multi-Point SHA-256 (Sampling 48KB) |
| **Performance** | Sub-second analysis on 91,646 assets |

---

## Project Status: 1.1 (Final)
We've achieved **Complete System Delivery** with a focus on minimalism. The tool has been stripped of unnecessary automation in favor of actionable, human-readable advice.
