# Purgomatic 💎

**The SRE Photo Librarian** — A high-performance, 100% Pure Go tool for auditing and planning large photo/video migrations (90k+ files).

---

## Philosophy: Home-First Intelligence
Purgomatic is built on the **"Home-First"** principle. Instead of just finding duplicates, it identifies **"Golden Winners"** (files already safely archived in your target Synology/Archive folders) and helps you eliminate the "Toil" (redundant copies on Phones/Thumb drives).

### Key Features
- **100% Pure Go (Zero CGO)**: No external C-dependencies. Built with `modernc.org/sqlite`.
- **Multi-Point Hashing**: Custom SHA-256 sampler (First/Middle/End) for collision-resistant deduplication at scale.
- **Stat-First Sync**: Blazing fast rescans by skipping hashing for unchanged files based on Size/Mtime.
- **SRE Insight Dashboard**: High-level strategic report on library health, toil, and historical density.
- **Home-Aware Planning**: Generates a 91k+ line `migration.json` for precise record keeping.

---

## Getting Started

### 1. Define Scan Targets
Create `scans.json` to tell Purgomatic where your "Source" and "Home" (Target) folders are:
```json
[
  { "path": "/Volumes/Synology/Photos" },
  { "path": "/Users/sam/Pictures/Imports" },
  { "path": "/Users/sam/iPhone_Backup" }
]
```

### 2. Standard Workflow
Using the provided `Makefile`:
```bash
make scan    # Index metadata (Parallel Hashing + Transactional SQLite)
make report  # Global SRE Insight & Strategic Advice
make plan    # Generate Home-First migration.json
```

---

## Architectural Specs
| Metric | Specification |
| :--- | :--- |
| **Language** | Go 1.23+ (Pure Go) |
| **Database** | SQLite (CGO-free) |
| **Concurrency** | 20 Worker Pool |
| **Hashing** | Multi-Point SHA-256 (Sampling 48KB) |
| **Performance** | Sub-second analysis on 91,646 assets |
