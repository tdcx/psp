#!/usr/bin/env bash
set -euo pipefail

FILE="docs/backlog/labels.yml"

if ! command -v yq &>/dev/null; then
  echo "❌ Necesitas instalar yq (brew install yq)"
  exit 1
fi
if ! command -v jq &>/dev/null; then
  echo "❌ Necesitas instalar jq (brew install jq)"
  exit 1
fi
if ! command -v gh &>/dev/null; then
  echo "❌ Necesitas instalar GitHub CLI (brew install gh)"
  exit 1
fi

echo "📂 Importando labels desde $FILE en el repo actual…"
count=0

yq -o=json '.[]' "$FILE" | jq -c '.' | while read -r row; do
  name=$(echo "$row" | jq -r '.name')
  color=$(echo "$row" | jq -r '.color')
  desc=$(echo "$row" | jq -r '.description')

  echo "➡️  Procesando label: $name"
  if gh label create "$name" --color "$color" --description "$desc" 2>/dev/null; then
    echo "   ✅ Creado"
  else
    gh label edit "$name" --color "$color" --description "$desc"
    echo "   🔄 Actualizado"
  fi
  count=$((count+1))
done

echo "🎉 Labels procesados: $count"
