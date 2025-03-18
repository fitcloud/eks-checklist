---
layout: default
title: EKS Checklist Docs
---

# EKS Checklist Docs

이 문서는 EKS 클러스터 점검을 위한 런북을 제공합니다.

## 📌 런북 목록

{% assign categories = site.pages | group_by_exp: "page", "page.dir | split: '/' | second" %}
{% for category in categories %}
### {{ category.name | capitalize }}

{% for page in category.items %}
- [{{ page.title }}]({{ page.url | relative_url }})
{% endfor %}

{% endfor %}
