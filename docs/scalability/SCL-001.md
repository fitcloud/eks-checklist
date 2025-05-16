# SCL-001 Karpenter 사용

## Meaning
Karpenter는 Kubernetes 클러스터에서 동적으로 노드를 자동으로 프로비저닝하고, 리소스 요구 사항에 맞는 최적의 노드 크기를 선택하여 클러스터의 리소스 효율성을 높이는 오픈 소스 자동화 툴입니다. Karpenter는 클러스터의 스케일링을 관리하고, 리소스가 부족하거나 오버프로비저닝일 때 이를 자동으로 조정합니다. 이는 특히 클라우드 환경에서 비용 최적화와 리소스 활용도를 개선하는 데 유용합니다.

## Impact
- 리소스 최적화: Karpenter는 애플리케이션의 리소스 요구사항에 따라 자동으로 노드를 추가하거나 삭제하여 클러스터의 리소스를 최적화합니다. 이로 인해, 필요 이상으로 많은 자원이 소비되는 상황을 방지하고, 불필요한 비용을 절감할 수 있습니다.
- 자동화된 스케일링: Karpenter는 트래픽 증가나 리소스 요구사항에 따라 자동으로 스케일링을 진행하여, 사용자가 수동으로 노드를 관리할 필요 없이 클러스터가 안정적으로 운영될 수 있도록 지원합니다. 이는 클러스터의 확장성과 탄력성을 크게 향상시킵니다.
- 비용 절감: Karpenter는 클라우드 환경에서 비용을 최적화할 수 있게 도와줍니다. 불필요한 노드를 줄이고 다양한 인스턴스 타입을 활용하여, 필요할 때만 리소스를 동적으로 할당함으로써 클라우드 서비스 사용 비용을 절감할 수 있습니다.
- 운영 효율성 증가: Karpenter의 자동화된 프로비저닝과 스케일링 기능은 클러스터 관리의 복잡도를 줄이고, 운영 효율성을 크게 향상시킵니다. 클러스터를 관리하는 데 드는 시간과 노력을 절감할 수 있습니다.

## Diagnosis
Karpenter가 설치되어 있는지 확인하세요

```bash
kubectl get deployments -n karpenter | grep "karpenter" >/dev/null && echo "Karpenter is installed" || echo "Karpenter is not installed"
```

## Mitigation
Karpenter를 설치하세요

example
```bash
helm upgrade --install karpenter oci://public.ecr.aws/karpenter/karpenter --version "${KARPENTER_VERSION}" --namespace "${KARPENTER_NAMESPACE}" --create-namespace \
  --set "settings.clusterName=${CLUSTER_NAME}" \
  --set "settings.interruptionQueue=${CLUSTER_NAME}" \
  --set controller.resources.requests.cpu=1 \
  --set controller.resources.requests.memory=1Gi \
  --set controller.resources.limits.cpu=1 \
  --set controller.resources.limits.memory=1Gi \
  --wait
```
[Karpenter 설치](https://karpenter.sh/docs/getting-started/getting-started-with-karpenter/#4-install-karpenter)
[AWS 모범사례 Karpenter](https://docs.aws.amazon.com/ko_kr/eks/latest/best-practices/karpenter.html)

