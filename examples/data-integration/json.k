import json

schema Server:
    ports: [int]

server: Server = json.decode('{"ports": [80, 8080]}')
server_json = json.encode({
    ports = [80, 8080]
})
