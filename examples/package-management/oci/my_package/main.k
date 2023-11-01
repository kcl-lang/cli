# Import and use the contents of the external dependency 'k8s'.
import k8s.api.core.v1 as k8core

k8core.Pod {
    metadata.name = "web-app"
    spec.containers = [{
        name = "main-container"
        image = "nginx"
        ports = [{containerPort = 80}]
    }]
}
