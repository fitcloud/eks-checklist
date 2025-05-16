# NET-002 Pod에 부여할 IP 부족시 알림 설정

## Meaning
Kubernetes 클러스터에서 Pod에 IP를 부여할 수 없을 때는 네트워크 리소스 부족을 나타냅니다. IP 부족은 Pod가 실행되지 않거나 네트워크 연결이 불안정해지는 문제를 유발할 수 있습니다. 이러한 상황을 알림 시스템을 통해 빠르게 감지하고 대응할 수 있도록 설정합니다.

## Impact
- Pod 스케줄링 실패: IP 할당 부족으로 새로운 Pod가 스케줄되지 않거나 실행되지 않을 수 있음
- 애플리케이션 장애: 네트워크 리소스 부족으로 애플리케이션의 성능 저하 또는 장애가 발생할 수 있음

## Diagnosis
Pod IP 부족 상태를 모니터링하려면, Kubernetes 이벤트를 통해 IP 할당 상태를 확인합니다.

```bash
kubectl get events --sort-by='.lastTimestamp' | grep 'failed to allocate a network IP'
```

## Mitigation
**eks addones으로 vpc-cni을 사용 중이라면 선택할 수 있습니다**

vpc-cni 배포 시 **`cni-metrics-helper`** 플러그인을 활성화 하여 추가 메트릭을 노출할 수 있습니다

플러그인의 경우 배포 시  매계변수`cniMetricsHelper.enabled=true` 를 활성화 함으로 사용 가능합니다

추가 메트릭 중  `totalIPAddresses` 항목을 통해 Pod의 할당된 IP 개수을 확인 가능합니다

해당 메트릭을 Prometheus & CloudWatch와 같은 서비스에서 집계하여 알람을 구성할 수 있습니다

이렇게 하면 사전에 IP 부족에 대해 인지하고 방지할 수 있습니다

**VPC IPAMs 메트릭 활용**

Amazon VPC IP Address Manager > IPAMs을 활성화 하고 사용하면 CloudWatch을 통해 Subnet IP 사용현황에 대해 Metrics을 추가하고 관리할 수 있습니다

메트릭 **`TotalActiveIpCount`** 은 활성화된 IP 개수를 의미하며 IP 개수 대비 비율 계산이 필요하나 CloudWatch Alarm을 구성하면 사전에 IP 부족에 대해 예방할 수 있습니다

[CNI-Metric-Helper](https://github.com/aws/amazon-vpc-cni-k8s/tree/master/cmd/cni-metrics-helper)
[Cloudwatch Ipam](https://docs.aws.amazon.com/ko_kr/vpc/latest/ipam/cloudwatch-ipam.html)