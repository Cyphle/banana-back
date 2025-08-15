# Kubernetes Configuration for Banana Backend

## Prerequisites
- Minikube running
- kubectl configured
- NGINX Ingress Controller enabled in minikube

## Structure
```
k8s/
├── postgres/
│   ├── postgres-deployment.yaml
│   └── postgres-service.yaml
├── keycloak/
│   ├── keycloak-deployment.yaml
│   ├── keycloak-service.yaml
│   ├── keycloak-config.yaml
│   ├── kustomization.yaml
│   ├── generate-keycloak-config.sh
│   └── README.md
├── redis/
│   ├── redis-deployment.yaml
│   └── redis-service.yaml
├── ingress.yaml
├── Makefile
└── README.md
```

## Configuration Keycloak - Source unique

⚠️ **Important** : La configuration Keycloak est maintenue dans un seul fichier source :
```
keycloak/realms/banana.json  ← Source unique de vérité
```

Trois solutions sont disponibles pour synchroniser cette configuration vers Kubernetes :

### 🚀 Solution 1 : Script de génération (Simple)

```bash
cd k8s/keycloak
./generate-keycloak-config.sh  # Génère keycloak-config.yaml
kubectl apply -f keycloak-config.yaml
```

### 🔧 Solution 2 : Kustomize (Natif K8s)

```bash
cd k8s/keycloak
kubectl apply -k .  # Utilise kustomization.yaml
```

**Kustomize** génère automatiquement le ConfigMap depuis le fichier source sans duplication.

### ⚡ Solution 3 : Makefile (Recommandé - Automatisé)

```bash
cd k8s
make deploy-keycloak  # Synchronise et déploie automatiquement
# Ou
make deploy-all       # Déploie tous les services
```

Le Makefile automatise complètement le workflow de synchronisation.

## Setup

1. Enable NGINX Ingress Controller:
```bash
minikube addons enable ingress
```

2. Deploy all services:

**Option A : Déploiement traditionnel**
```bash
# Deploy individual services
kubectl apply -f postgres/
kubectl apply -f keycloak/keycloak-config.yaml  # Généré au préalable
kubectl apply -f keycloak/keycloak-deployment.yaml
kubectl apply -f keycloak/keycloak-service.yaml
kubectl apply -f redis/
kubectl apply -f ingress.yaml
```

**Option B : Avec Kustomize (Recommandé)**
```bash
kubectl apply -k keycloak/  # Génère automatiquement la config
kubectl apply -f postgres/
kubectl apply -f redis/
kubectl apply -f ingress.yaml
```

**Option C : Avec Makefile (Plus simple)**
```bash
make deploy-all  # Synchronise et déploie tout automatiquement
```

## Access Services

### 🔌 NodePort - Accès direct (Recommandé pour le développement)

Les **NodePort** permettent d'accéder aux services directement depuis votre machine :

- **PostgreSQL**: `localhost:30432` - Base de données
- **Keycloak**: `localhost:30181` - Interface d'authentification  
- **Redis**: `localhost:30379` - Cache

**Pourquoi NodePort ?**
- ✅ **Accès direct** : Pas besoin d'Ingress ou de port-forward
- ✅ **Simple** : Fonctionne immédiatement avec minikube
- ✅ **Développement** : Idéal pour tester et déboguer
- ✅ **Stabilité** : Ports fixes, pas de changement

**Comment ça marche ?**
```bash
# Minikube expose automatiquement les NodePort sur localhost
minikube service banana-postgresql-service --url  # Affiche l'URL complète
minikube service keycloak-service --url
minikube service redis-service --url
```

### 🌐 Ingress - Accès via nom de domaine (Production)

L'**Ingress** permet d'accéder via des URLs avec noms de domaine :

```bash
# Activer le tunnel minikube (requis pour Ingress)
minikube tunnel
```

Puis accéder via :
- **Keycloak**: `http://localhost/keycloak`
- **PostgreSQL**: `http://localhost/postgres` 
- **Redis**: `http://localhost/redis`

**Différence NodePort vs Ingress :**
- **NodePort** = Accès direct par port (`:30181`)
- **Ingress** = Accès par chemin URL (`/keycloak`)

## Port Mapping from Docker Compose
- PostgreSQL: 5434 → 30432
- Keycloak: 8181 → 30181
- Redis: 6379 → 30379

## Commands
```bash
# Get minikube IP
minikube ip

# Get service URLs
minikube service banana-postgresql-service --url
minikube service keycloak-service --url
minikube service redis-service --url

# Clean up
kubectl delete -f k8s/
```