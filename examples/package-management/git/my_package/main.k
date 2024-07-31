import konfig.models.kube.frontend

frontend.Server {
    name = "nginx-deployment"
    image = "nginx"
    mainContainer.name = name
}
