import k8s.api.core.v1
import regex

protocol PVCProtocol:
    pvc?: {str:}

mixin PersistenVolumeClaimMixin for PVCProtocol:
    """
    PersistentVolumeClaim (PVC) sample:
    Link: https://kubernetes.io/docs/concepts/storage/persistent-volumes/#persistentvolumeclaims
    """

    # Mix in a new attribute `kubernetesPVC`
    kubernetesPVC: v1.PersistentVolumeClaim

    if pvc:
        kubernetesPVC = v1.PersistentVolumeClaim {
            metadata.name = pvc.name
            metadata.labels = pvc.labels
            spec = {
                accessModes = pvc.accessModes
                resources = pvc.resources
                storageClassName = pvc.storageClassName
                }
        }

schema Server(Deployment):
    mixin [PersistenVolumeClaimMixin]
    pvc?: {str:}
    """ pvc user interface data defined by PersistentVolumeClaimMixin """

schema Deployment[priority: int]:
    name: str
    cpu: int = _cpu
    memory: int = _cpu * 2
    volumes?: [Volume]
    image: str
    service?: Service
    replica: int = 1
    command: [str]
    labels?: {str: str}

    if priority == 1:
        _cpu = 256
    elif priority == 2:
        _cpu = 512
    elif priority == 3:
        _cpu = 1024
    else:
        _cpu = 2048

    check:
        multiplyof(cpu, 256), "CPU must be a multiple of 256"
        regex.match(image, "^[a-zA-Z]+:\d+\.\d+\.\d+$"), "image name should be like 'nginx:1.14.2'"
        1 <= replica <= 100, "replica must be between 1 and 100"
        len(labels) >= 2, "The length of labels should be greater than or equal to 2"
        "env" in labels, "'env' must be in labels"
        len(command) > 0, "command-list should be non-empty"

schema Port:
    name: str
    protocol: str
    port: int
    targetPort: int

    check:
        port in [80, 443], "port must be 80 or 443"
        protocol in ["TCP", "HTTP"], "protocol must be TCP or HTTP"
        1024 < targetPort, "targetPort must be greater than 1024"

schema Service:
    name: "my-service" = "my-service"
    ports: [Port]

    check:
        len(ports) > 0, "ports-list should be non-empty"

schema Volume:
    name: str
    mountPath: str
    hostPath: str

server = Server {
    name = "my-nginx"
    image = "nginx:1.14.2"
    volumes = [Volume {
        name = "mydir"
        mountPath = "/test-pd"
        hostPath = "/data"
    }]
    command = ["nginx"]
    labels = {
        run = "my-nginx"
        env = "pre-prod"
    }
    service.ports = [Port {
        name = "http"
        protocol = "TCP"
        port = 80
        targetPort = 9376
    }]
    pvc = {
        name = "my_pvc"
        acccessModes = ["ReadWriteOnce"]
        resources.requests.storage = "8Gi"
        storageClassName = "slow"
    }
}
