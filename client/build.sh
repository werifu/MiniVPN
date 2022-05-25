go build
docker cp ./client HostU:/
docker cp ./client.key HostU:/
docker cp ./client.crt HostU:/
docker cp ./config.json HostU:/
docker exec -it HostU /bin/bash