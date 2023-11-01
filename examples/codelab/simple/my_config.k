_priority = option("priority") # a non-exported and mutable variable
_env = option("env") # a non-exported and mutable variable
_priorityCpuMap = {
    "1" = 256
    "2" = 512
    "3" = 1024
}
# Using a dict to simplify logic and the default value is 2048
_cpu = _priorityCpuMap[_priority] or 2048
_name = "nginx"
# exported variables
cpu = _cpu
memory = _cpu * 2
command = [_name] # a list
labels = {
    run = "my-{}".format(_name)
    env = _env
} # a dict
image = "{}:1.14.2".format(_name) # string format
service = "my-service"
