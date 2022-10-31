FROM hyperledger/fabric-tools:2.4.6

WORKDIR /app
COPY ./ /hlf-operator

RUN go build -o kubectl-hlf /app/kubectl-hlf/main.go

