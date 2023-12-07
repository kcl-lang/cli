run:
	kcl import -m crd crd.yaml
	cd models && rm -r ./k8s
	rm -rf ./models/dapr_io_v1alpha1_subscription.k
	kcl mod init
	kcl mod add k8s
	kcl run ./models
	cd ..

test:
	make run
