import regex

schema Deployment[priority]:
    name: str
    volumes?: [Volume]
    image: str
    service?: Service
    replica: int = 1
    command: [str]
    labels?: {str:str}

    if priority == 1:
        _cpu = 256
    elif priority == 2:
        _cpu = 512
    elif priority == 3:
        _cpu = 1024
    else:
        _cpu = 2048

    cpu: int = _cpu
    memory: int = _cpu * 2

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

nginx = Deployment(priority=2) {
    name = "my-nginx"
    image = "nginx:1.14.2"
    volumes = [Volume {
        name = "mydir"
        mountPath = "/test-pd"
        hostPath = "/data"
    }]
    command = ["nginx"]
    labels.run = "my-nginx"
    labels.env = "pre-prod"
    service.ports = [Port {
        name = "http"
        protocol = "TCP"
        port = 80
        targetPort = 9376
    }]
}
