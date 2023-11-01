import regex
import k8spkg.api.core.v1

protocol PVCProtocol:
    pvc?: {str:}

mixin PersistentVolumeClaimMixin for PVCProtocol:
    """
    PersistentVolumeClaim (PVC) sample:
    Link: https://kubernetes.io/docs/concepts/storage/persistent-volumes/#persistentvolumeclaims
    """
    
    # Mix in a new attribute `kubernetesPVC`
    kubernetesPVC?: v1.PersistentVolumeClaim

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

schema Deployment:
    name: str
    cpu: int = 512
    memory: int = cpu * 2
    volumes?: [Volume]
    image: str
    service?: Service
    replica: int = 1
    command: [str]
    labels?: {str:str}

    check:
        multiplyof(cpu, 256), "cpu must be a multiplier of 256"
        regex.match(image, "^[a-zA-Z]+:\d+\.\d+\.\d+$"), "image name should be like 'nginx:1.14.2'"
        1 <= replica < 100, "replica should be in range (1, 100)"
        len(labels) >= 2 if labels, "the length of labels should be large or equal to 2"
        "env" in labels, "'env' must be in labels"
        len(command) > 0, "the command list should be non-empty"

schema Port:
    name: str
    protocol: str
    port: int
    targetPort: int

schema Service:
    name: "my-service" = "my-service"
    ports: [Port]

schema Volume:
    name: str
    mountPath: str
    hostPath: str

schema Server(Deployment):
    mixin [PersistentVolumeClaimMixin]
    pvc?: {str:}
    """ pvc user interface data defined by PersistentVolumeClaimMixin """

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
        accessModes = ["ReadWriteOnce"]
        resources.requests.storage = "8Gi"
        storageClassName = "slow"
    }
}
