
To build the image for nem-onmand-service use following command.

docker build -t nem-ondemand-proxy .


To launch the docker cotainer for nem-ondemand-proxy use following command

docker-compose -f compose/nem-proxy-go.yml up -d


To submit a request:

# Enter the nem-ondemand-api container

# install curl and download grpcurl
apk add curl && curl -L https://github.com/fullstorydev/grpcurl/releases/download/v1.4.0/grpcurl_1.4.0_linux_x86_64.tar.gz | tar -xz && mv grpcurl /usr/bin/

# execute a request on the first bbsim onu
grpcurl -plaintext -d '{"id": "2910b26bbb29521d93fab21b"}' localhost:50052 on_demand_api.NemService/OmciTest
