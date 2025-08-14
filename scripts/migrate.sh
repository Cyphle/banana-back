#!/bin/bash

# Script de migration utilisant la configuration YAML
# Usage: ./scripts/migrate.sh [up|down|status]

set -e

# Vérifier que le script est exécuté depuis la racine du projet
if [ ! -f "Cargo.toml" ]; then
    echo "Erreur: Ce script doit être exécuté depuis la racine du projet"
    exit 1
fi

# Fonction pour extraire la DATABASE_URL depuis le YAML
get_database_url() {
    # Utiliser yq si disponible, sinon grep (solution basique)
    if command -v yq &> /dev/null; then
        # Avec yq (plus robuste)
        yq eval '.database | "postgres://\(.username):\(.password)@\(.host):\(.port)/\(.schema)"' config/default.yaml
    else
        # Solution basique avec grep (moins robuste)
        host=$(grep "host:" config/default.yaml | head -1 | awk '{print $2}' | tr -d '"')
        port=$(grep "port:" config/default.yaml | head -1 | awk '{print $2}' | tr -d '"')
        schema=$(grep "schema:" config/default.yaml | head -1 | awk '{print $2}' | tr -d '"')
        username=$(grep "username:" config/default.yaml | head -1 | awk '{print $2}' | tr -d '"')
        password=$(grep "password:" config/default.yaml | head -1 | awk '{print $2}' | tr -d '"')
        echo "postgres://${username}:${password}@${host}:${port}/${schema}"
    fi
}

# Obtenir la DATABASE_URL
DATABASE_URL=$(get_database_url)

if [ -z "$DATABASE_URL" ]; then
    echo "Erreur: Impossible de lire la configuration de la base de données"
    exit 1
fi

echo "Configuration de la base de données:"
echo "  Host: $(echo $DATABASE_URL | sed 's/.*@\([^:]*\):.*/\1/')"
echo "  Port: $(echo $DATABASE_URL | sed 's/.*:\([0-9]*\)\/.*/\1/')"
echo "  Database: $(echo $DATABASE_URL | sed 's/.*\///')"
echo ""

# Exporter la variable d'environnement
export DATABASE_URL

# Exécuter la commande de migration
case "${1:-up}" in
    "up")
        echo "🚀 Exécution des migrations..."
        sea-orm-cli migrate up
        ;;
    "down")
        echo "⬇️  Annulation des migrations..."
        sea-orm-cli migrate down
        ;;
    "status")
        echo "📊 Statut des migrations..."
        sea-orm-cli migrate status
        ;;
    *)
        echo "Usage: $0 [up|down|status]"
        echo "  up     - Exécuter les migrations"
        echo "  down   - Annuler les migrations"
        echo "  status - Afficher le statut des migrations"
        exit 1
        ;;
esac
