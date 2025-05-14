FROM golang:1.24-alpine AS builder

# 필수 패키지 설치
RUN apk add --no-cache git ca-certificates curl

# 작업 디렉토리 설정
WORKDIR /

# 소스 코드 복사
COPY . .

# 의존성 다운로드 및 빌드
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o eks-checklist .

# 최종 이미지 생성
FROM alpine:3.18

# 필수 패키지 설치
RUN apk add --no-cache python3 py3-pip curl unzip bash jq findutils

# AWS CLI 설치
RUN pip3 install --upgrade pip && \
    pip3 install awscli

# kubectl 설치
RUN curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl" && \
    chmod +x kubectl && \
    mv kubectl /usr/local/bin/

# 작업 디렉토리 설정
WORKDIR /

# 애플리케이션 복사
COPY --from=builder /eks-checklist /
COPY --from=builder /templates /templates

# 필요한 디렉토리 생성
RUN mkdir -p /root/.kube /root/.aws /output

# # 스크립트 추가
# COPY entrypoint.sh /
# RUN chmod +x /entrypoint.sh

# 환경 변수 설정
ENV KUBECONFIG=/root/.kube/config
ENV AWS_SHARED_CREDENTIALS_FILE=/root/.aws/credentials
ENV AWS_CONFIG_FILE=/root/.aws/config
ENV RUNNING_IN_DOCKER=true

# 작업 디렉토리 유지 (원본 애플리케이션에 영향 없도록)
WORKDIR /

ENTRYPOINT ["./eks-checklist"]
# CMD [""]