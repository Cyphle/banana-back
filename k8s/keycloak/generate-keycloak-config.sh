#!/bin/bash

# Script pour générer automatiquement le ConfigMap Keycloak depuis le fichier source
# Usage: cd k8s/keycloak && ./generate-keycloak-config.sh

set -e

SOURCE_FILE="../../keycloak/realms/banana.json"
OUTPUT_FILE="keycloak-config.yaml"

if [ ! -f "$SOURCE_FILE" ]; then
    echo "Erreur: Le fichier source $SOURCE_FILE n'existe pas"
    exit 1
fi

echo "Génération du ConfigMap Keycloak depuis $SOURCE_FILE..."

# Créer le début du fichier YAML
cat > "$OUTPUT_FILE" << 'EOF'
apiVersion: v1
kind: ConfigMap
metadata:
  name: keycloak-realm-config
data:
  banana.json: |
EOF

# Ajouter le contenu JSON avec indentation de 4 espaces
sed 's/^/    /' "$SOURCE_FILE" >> "$OUTPUT_FILE"

echo "✅ ConfigMap généré dans $OUTPUT_FILE"
echo "💡 Pour appliquer les changements:"
echo "   kubectl apply -f $OUTPUT_FILE"