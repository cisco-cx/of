from debian:buster

run apt-get update && apt-get install -y python3-pytest python3-yaml python3-requests && rm -rf /var/lib/apt

workdir /src/smoke_tests
