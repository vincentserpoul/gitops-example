.PHONY: deploy-latest context-select image-tag image-push manifest-kustomize-build manifests-apply release-commit

ENV?=dev

include env/base/base.env
include env/overlays/$(ENV)/$(ENV).env

export

IMAGE_NAME = $(CONTAINER_REGISTRY_URL)$(APP_NAME):$(APP_VERSION)
APP_TLS_SECRET = $(shell echo $(APP_DOMAIN) | tr "." "-")

deploy-latest: context-select image-tag image-push manifest-kustomize-build manifest-apply

# deploy-latest: context-select image-tag image-push manifest-kustomize-build manifest-apply release-commit


###########
# HELPERS #
###########

release-commit:
	git commit -m "manifest build and deploy $(IMAGE_NAME) to $(ENV)" ./manifests/*

context-select:
	kubectl config use-context $(K8S_CONTEXT)

############
# MANIFEST #
############

MANIFEST_FOLDER = ./manifests/$(ENV)
MANIFEST_FILE_PATH = $(MANIFEST_FOLDER)/$(PROJECT_NAME)_$(APP_NAME).yaml

manifest-apply:
	kubectl apply -f $(MANIFEST_FILE_PATH)

manifest-kustomize-build:
	kubectl kustomize ./env/overlays/$(ENV) | \
	sed \
	-e "s/APP_VERSION/$(APP_VERSION)/g" \
	-e "s/APP_PORT/$(APP_PORT)/g" \
	-e "s/APP_DOMAIN/$(APP_DOMAIN)/g" \
	-e "s/APP_TLS_SECRET/$(APP_TLS_SECRET)/g" \
	-e "s/APP_NAME/$(APP_NAME)/g" \
	-e "s/PROJECT_NAME/$(PROJECT_NAME)/g" \
	-e "s/CONTAINER_REGISTRY_URL/$(CONTAINER_REGISTRY_URL)/g" \
	-e "s/VAR_APP_ENV/$(ENV)/g" \
	> $(MANIFEST_FILE_PATH)

#########
# IMAGE #
#########

image-tag:
	docker tag $(APP_NAME) $(IMAGE_NAME)

image-push:
	@if [ -z "$(CONTAINER_REGISTRY_URL)" ]; then \
		echo "no container registry to push to as URL is empty"; \
	else \
		docker push $(IMAGE_NAME); \
	fi
