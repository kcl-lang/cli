import manifests

import app

# Convert the `App` model into Kubernetes Deployment and Service Manifests
test_kubernetesRender = lambda {
    a = app.App {
        name = "app"
        containers.ngnix = {
            image = "ngnix"
            ports = [{containerPort = 80}]
        }
        service.ports = [{ port = 80 }]
    }
    deployment_got = kubernetesRender(a)
    assert deployment_got[0].kind == "Deployment"
    assert deployment_got[1].kind == "Service"
}
