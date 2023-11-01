import pkg

server: pkg.Server {
    # Override the base image.
    image = "nginx:1.14.2-dev"
    # Union a new label env into base labels.
    labels.env = "dev"
    # Append a port into base ports.
    mainContainer.ports += [{
        protocol = "TCP"
        port = 443
        targetPort = 1100
    }]
}
