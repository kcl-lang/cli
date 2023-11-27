import pkg.app as app_pkg
import manifests

app: app_pkg.App {
    domainType = "Standard"
    containerPort = 80
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
render = lambda app: app_pkg.App {
    {
        apiVersion = "v1"
        metadata.name = app.name
    }
}

manifests.yaml_stream([render(a) for a in app_pkg.App.instances()])
