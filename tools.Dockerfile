FROM hyperledger/fabric-tools:2.4.6

WORKDIR /app
COPY ./ /app
RUN go build -o kubectl-hlf-bin /app/kubectl-hlf/main.go
RUN chmod +x ./kubectl-hlf-bin
RUN mv ./kubectl-hlf-bin /usr/local/bin/kubectl-hlf
RUN apk add curl
RUN curl -LO https://storage.googleapis.com/kubernetes-release/release/$(curl -s https://storage.googleapis.com/kubernetes-release/release/stable.txt)/bin/linux/amd64/kubectl
RUN chmod +x ./kubectl
RUN mv ./kubectl /usr/local/bin
