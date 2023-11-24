docker:
	kcl main.k docker_compose_render.k

k8s:
	kcl main.k kubernetes_render.k

kcl-test:
	kcl test

test:
	make docker
	make k8s
	make kcl-test
