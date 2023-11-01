import yaml

schema Server:
    ports: [int]

server: Server = yaml.decode("""\
ports:
- 80
- 8080
""")
server_yaml = yaml.encode({
    ports = [80, 8080]
})
