docker:
	kcl main.k docker_compose_render.k

k8s:
	kcl main.k kubernetes_render.k

test:
	make docker
	make k8s