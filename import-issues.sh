#!/usr/bin/env bash
set -euo pipefail

CSV="docs/backlog/github_backlog.csv"
REPO="$(gh repo view --json nameWithOwner -q .nameWithOwner)"

if ! command -v gh >/dev/null; then echo "❌ Falta GitHub CLI (gh)"; exit 1; fi
if ! command -v python3 >/dev/null; then echo "❌ Falta Python3"; exit 1; fi

python3 - <<'PY'
import csv, json, subprocess, sys, os

csv_path = "docs/backlog/github_backlog-2.csv"

# Carga milestones existentes (nombre -> número)
repo = subprocess.check_output(["gh","repo","view","--json","nameWithOwner","-q",".nameWithOwner"]).decode().strip()
milestones = json.loads(subprocess.check_output(["gh","api",f"/repos/{repo}/milestones"]).decode())
ms_map = { m["title"]: m["number"] for m in milestones }

def ensure_milestone(title):
    if not title: return None
    if title in ms_map: return ms_map[title]
    # crea si no existe
    created = json.loads(subprocess.check_output([
        "gh","api","-X","POST",f"/repos/{repo}/milestones","-f",f"title={title}"
    ]).decode())
    ms_map[title] = created["number"]
    return created["number"]

with open(csv_path, newline='', encoding='utf-8') as f:
    reader = csv.DictReader(f)
    for i,row in enumerate(reader, start=1):
        title = row.get("Title","").strip()
        body  = row.get("Body","").strip()
        labels = [l.strip() for l in row.get("Labels","").split(",") if l.strip()]
        milestone_title = row.get("Milestone","").strip()
        ms_number = ensure_milestone(milestone_title)

        payload = {"title": title, "body": body}
        if labels: payload["labels"] = labels
        if ms_number: payload["milestone"] = ms_number

        print(f"#{i} creando issue: {title}")
        # gh api para crear la issue
        subprocess.check_call([
            "gh","api","-X","POST", f"/repos/{repo}/issues",
            "-H","Accept: application/vnd.github+json",
            "-f", f"title={title}",
            "-f", f"body={body}",
            *sum([["-f", f"labels[]={lbl}"] for lbl in labels], []),
            *(["-f", f"milestone={ms_number}"] if ms_number else [])
        ])
PY
