import pkg

server = pkg.Server {
    name = "my-nginx"
    image = "nginx:1.14.2"
    volumes = [pkg.Volume {
        name = "mydir"
        mountPath = "/test-pd"
        hostPath = "/data"
    }]
    command = ["nginx"]
    labels = {
        run = "my-nginx"
        env = "pre-prod"
    }
    service.ports = [pkg.Port {
        name = "http"
        protocol = "TCP"
        port = 80
        targetPort = 9376
    }]
}
