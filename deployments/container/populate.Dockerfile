FROM ubuntu:20.04

WORKDIR /src

RUN apt-get update && apt-get install -y curl unzip

RUN curl "https://awscli.amazonaws.com/awscli-exe-linux-x86_64.zip" -o "awscli2.zip"
RUN unzip awscli2.zip
RUN ./aws/install

COPY . .

ENTRYPOINT [ "/src/scripts/populate_entrypoint.sh" ]
