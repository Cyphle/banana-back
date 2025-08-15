# Configuration Keycloak pour Kubernetes

## Solutions pour éviter la duplication de configuration

Voici 3 solutions pour maintenir une seule source de vérité pour la configuration Keycloak :

### 🚀 Solution 1 : Script de génération (Recommandée)

**Usage :**
```bash
cd k8s/keycloak
./generate-keycloak-config.sh
kubectl apply -f keycloak-config.yaml
```

**Avantages :**
- Simple et direct
- Génération automatique du ConfigMap
- Pas de dépendances externes

### 🔧 Solution 2 : Kustomize

**Usage :**
```bash
cd k8s/keycloak
kubectl apply -k .
```

**Avantages :**
- Outil natif Kubernetes
- Référence directe au fichier source
- Pas de duplication

### ⚡ Solution 3 : Makefile (Plus pratique)

**Usage :**
```bash
cd k8s
make deploy-keycloak  # Synchronise et déploie automatiquement
```

**Avantages :**
- Automatisation complète
- Une seule commande pour tout
- Workflow cohérent

## Recommandation

Pour un workflow optimal, utilisez **Solution 3 (Makefile)** qui combine simplicité et automation :

```bash
# Déploie Keycloak avec config à jour
make deploy-keycloak

# Ou déploie tout
make deploy-all
```

## Fichier source unique

La configuration Keycloak est maintenue dans :
```
keycloak/realms/banana.json  ← Source unique de vérité
```

Toutes les solutions référencent automatiquement ce fichier.