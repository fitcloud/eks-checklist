# REL-012 Cluster Autoscaler 적용

## Meaning
Cluster Autoscaler는 클라우드 환경에서 Kubernetes 클러스터의 노드를 자동으로 확장하거나 축소하여 워크로드에 필요한 자원을 동적으로 제공하는 도구입니다. 이를 통해 클러스터의 리소스를 효율적으로 관리하고 비용을 최적화할 수 있습니다

## Impact
- 리소스스 부족 :  워크로드의 수요가 증가해도 자동으로 노드가 추가되지 않으면, CPU, 메모리, 또는 기타 리소스가 부족해져 Pod가 Pending 상태로 남게 됩니다. 이는 애플리케이션 성능 저하나 장애를 유발할 수 있습니다
- 서비스 중단 :  충분한 리소스가 없으면 신규 요청을 처리하지 못하거나 기존 서비스가 중단될 수 있습니다.
- EKS Cluster 관리 복장성 증가 : Node의 수와 용량을 수동으로 관리해야 하므로, 클러스터 관리가 매우 복잡해지고 시간이 많이 소요됩니다. 클러스터의 부하나 리소스 수요 변화에 따라 수동으로 노드를 추가하거나 제거해야 하므로, 운영 부담이 커집니다.
- 비용 증가 : 필요 이상의 리소스를 할당받게 되어, 클라우드 비용이 불필요하게 커질 수 있습니다. Cluster Autoscaler가 없으면 자원을 최적화하기 어렵습니다.

## Diagnosis
ClusterAutoscaler가 설치되어 있는지 확인하세요

```bash
kubectl get deployment -n kube-system | grep cluster-autoscaler
kubectl get pods -n kube-system | grep cluster-autoscaler
```

## Mitigation
ClusterAutosacler를 설치하세요

**example**
```bash
helm repo add autoscaler https://kubernetes.github.io/autoscaler
helm repo update
helm install cluster-autoscaler autoscaler/cluster-autoscaler \
  --namespace kube-system \
  --set autoDiscovery.clusterName=<cluster-name> \
  --set awsRegion=<region> \
  --set rbac.create=true
```

[Cluster Autoscaler 설치](https://github.com/kubernetes/autoscaler/blob/master/cluster-autoscaler/cloudprovider/aws/README.md)
[Kubernetes Autoscaling](https://docs.aws.amazon.com/ko_kr/eks/latest/userguide/autoscaling.html)


