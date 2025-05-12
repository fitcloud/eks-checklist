#!/bin/bash
# filepath: /workspaces/eks-checklist/eks-check.sh
set -e

# 스크립트 디렉토리 확인
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# 현재 디렉토리 경로
CURRENT_DIR="$(pwd)"

# 도움말 출력 함수
show_help() {
  echo "EKS Checklist Docker 실행 스크립트"
  echo ""
  echo "사용법: $0 [옵션]"
  echo ""
  echo "옵션:"
  echo "  --help                   이 도움말 메시지를 표시합니다."
  echo "  --context CONTEXT        사용할 kubeconfig 컨텍스트 지정"
  echo "  --profile PROFILE        사용할 AWS 프로필 지정"
  echo "  --filter FILTER          출력 필터 (all, pass, fail, manual)"
  echo "  --output FORMAT          출력 형식 (text, html)"
  echo "  --sort                   결과를 상태별로 정렬"
  echo ""
  echo "예제:"
  echo "  $0 --context my-eks-cluster --output html"
  echo "  $0 --profile dev --filter fail"
}

# 기본 변수 설정
CONTEXT=""
PROFILE=""
FILTER=""
OUTPUT=""
SORT=""
RESULT_DIR="${CURRENT_DIR}"
INTERACTIVE=""

# 명령줄 인자 처리
while [[ $# -gt 0 ]]; do
  case $1 in
    --help)
      show_help
      exit 0
      ;;
    --context)
      CONTEXT="$2"
      shift 2
      ;;
    --profile)
      PROFILE="$2"
      shift 2
      ;;
    --filter)
      FILTER="$2"
      shift 2
      ;;
    --output)
      OUTPUT="$2"
      shift 2
      ;;
    --sort)
      SORT="--sort"
      shift
      ;;
    *)
      echo "알 수 없는 옵션: $1"
      show_help
      exit 1
      ;;
  esac
done

# 결과 디렉토리 생성
mkdir -p "${RESULT_DIR}"

# Docker 명령 구성
DOCKER_CMD="docker run --rm"

# 대화형 모드 설정
if [ -t 0 ] && [ -t 1 ] || [ -n "$INTERACTIVE" ]; then
  DOCKER_CMD="$DOCKER_CMD -it"
fi

# 볼륨 마운트 설정
DOCKER_CMD="$DOCKER_CMD -v ${HOME}/.kube:/root/.kube -v ${HOME}/.aws:/root/.aws -v ${RESULT_DIR}:/output"

# 환경 변수 설정
if [ -n "$CONTEXT" ]; then
  DOCKER_CMD="$DOCKER_CMD -e CONTEXT=$CONTEXT"
fi

# 이미지 및 추가 인자 설정
DOCKER_CMD="$DOCKER_CMD eks-checklist"

# 각 옵션 추가
if [ -n "$CONTEXT" ]; then
  DOCKER_CMD="$DOCKER_CMD --context $CONTEXT"
fi

if [ -n "$PROFILE" ]; then
  DOCKER_CMD="$DOCKER_CMD --profile $PROFILE"
fi

if [ -n "$FILTER" ]; then
  DOCKER_CMD="$DOCKER_CMD --filter $FILTER"
fi

if [ -n "$OUTPUT" ]; then
  DOCKER_CMD="$DOCKER_CMD --output $OUTPUT"
fi

if [ -n "$SORT" ]; then
  DOCKER_CMD="$DOCKER_CMD --sort"
fi

echo "명령 실행: $DOCKER_CMD"
eval "$DOCKER_CMD"

echo "결과물은 output 디렉토리에 저장되었습니다."