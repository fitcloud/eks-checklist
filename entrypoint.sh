#!/bin/bash
# filepath: /workspaces/eks-checklist/entrypoint.sh
set -e

echo "EKS Checklist 실행 중..."

# 결과 디렉토리 생성 및 권한 설정
mkdir -p /output
chmod 777 /output

# kubeconfig 파일 확인
if [ ! -f "$KUBECONFIG" ]; then
  echo "경고: kubeconfig 파일이 존재하지 않습니다."
  echo "호스트의 kubeconfig를 마운트하세요: -v ~/.kube:/root/.kube"
  exit 1
fi

# 결과물 저장 경로를 환경 변수로 설정 
export OUTPUT_DIR="/output"

# 인터랙티브 모드 감지
if [ -t 0 ] && [ -t 1 ]; then
  echo "대화형 모드로 실행 중입니다."
else
  echo "비대화형 모드로 실행 중입니다."
  
  # 비대화형 모드에서 컨텍스트가 명시적으로 지정되지 않은 경우
  if ! echo "$@" | grep -q -- "--context"; then
    # 현재 컨텍스트 사용
    CURRENT_CONTEXT=$(kubectl config current-context 2>/dev/null || echo "")
    if [ -n "$CURRENT_CONTEXT" ]; then
      echo "현재 컨텍스트를 사용합니다: $CURRENT_CONTEXT"
      set -- "--context" "$CURRENT_CONTEXT" "$@"
    else
      # 단일 컨텍스트 확인
      CONTEXTS_COUNT=$(kubectl config get-contexts -o name 2>/dev/null | wc -l)
      if [ "$CONTEXTS_COUNT" -eq 1 ]; then
        CONTEXT=$(kubectl config get-contexts -o name)
        echo "단일 컨텍스트를 자동으로 선택합니다: $CONTEXT"
        set -- "--context" "$CONTEXT" "$@"
      else
        echo "경고: 컨텍스트가 여러 개이거나 찾을 수 없습니다."
        echo "사용 가능한 컨텍스트:"
        kubectl config get-contexts -o name
      fi
    fi
  fi
fi

# EKS Checklist 실행
echo "실행 명령: /app/eks-checklist $@"
/app/eks-checklist "$@"

# 모든 결과물을 /output 디렉토리로 복사
echo "결과물을 /output 디렉토리로 복사 중..."

# 현재 디렉토리의 결과물 복사
if [ "$(pwd)" != "/output" ]; then
  # HTML 파일 복사
  find "$(pwd)" -maxdepth 1 -name "eks-checklist-report*.html" -exec cp -v {} /output/ \;
  
  # result 디렉토리가 있으면 복사
  if [ -d "$(pwd)/result" ]; then
    mkdir -p /output/result
    cp -rv "$(pwd)/result"/* /output/result/
  fi

  # results 디렉토리가 있으면 복사
  if [ -d "$(pwd)/results" ]; then
    mkdir -p /output/results
    cp -rv "$(pwd)/results"/* /output/results/
  fi
fi

# /app 디렉토리의 결과물 복사
if [ -d "/app/result" ]; then
  mkdir -p /output/result
  cp -rv /app/result/* /output/result/
fi

if [ -d "/app/results" ]; then
  mkdir -p /output/results
  cp -rv /app/results/* /output/results/
fi

find "/app" -maxdepth 1 -name "eks-checklist-report*.html" -exec cp -v {} /output/ \;

# 결과물 권한 설정
chmod -R 777 /output

echo "결과물이 /output 디렉토리에 저장되었습니다."