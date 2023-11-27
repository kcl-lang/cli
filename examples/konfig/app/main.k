import pkg.app as app_pkg

app: app_pkg.App {
    name: "app"
    domainType = "Standard"
    containerPort = 80
    if option("env") == "prod":
        volumes = [
            {
                mountPath = "/tmp"
            }
        ]
    services = [
        {
            clusterIP = "None"
            $type = "ClusterIP"
        }
    ]
}
