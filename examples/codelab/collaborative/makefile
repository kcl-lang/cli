base:
	kcl appops/test_app/base/base.k

dev:
	cd appops/test_app/dev && kcl -Y ./kcl.yaml

prod:
	cd appops/test_app/prod && kcl -Y ./kcl.yaml

default:
	make base
	make dev
	make prod
