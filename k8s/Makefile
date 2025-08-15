.PHONY: sync-keycloak-config deploy-keycloak deploy-all clean

# Synchronise la config Keycloak depuis le fichier source
sync-keycloak-config:
	@echo "🔄 Synchronisation de la configuration Keycloak..."
	cd keycloak && ./generate-keycloak-config.sh

# Déploie Keycloak avec la config à jour
deploy-keycloak: sync-keycloak-config
	@echo "🚀 Déploiement de Keycloak..."
	kubectl apply -f keycloak/

# Déploie tous les services
deploy-all:
	@echo "🚀 Déploiement de tous les services..."
	kubectl apply -R -f .

# Nettoie les ressources
clean:
	@echo "🧹 Nettoyage des ressources..."
	kubectl delete -R -f . --ignore-not-found=true

# Aide
help:
	@echo "Commandes disponibles:"
	@echo "  sync-keycloak-config  - Synchronise la config Keycloak"
	@echo "  deploy-keycloak       - Déploie Keycloak avec config à jour"
	@echo "  deploy-all           - Déploie tous les services"
	@echo "  clean                - Supprime toutes les ressources"