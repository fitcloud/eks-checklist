#!/bin/bash

set -e  # 에러 발생 시 스크립트 중단

# aws-runas 설치
arch=$(arch | sed s/aarch64/arm64/ | sed s/x86_64/amd64/) && \
wget https://github.com/mmmorris1975/aws-runas/releases/download/3.5.2/aws-runas-3.5.2-linux-${arch}.zip && \
unzip aws-runas-3.5.2-linux-${arch}.zip && \
chmod +x aws-runas && \
mv aws-runas /usr/local/bin/ && \
rm aws-runas-3.5.2-linux-${arch}.zip