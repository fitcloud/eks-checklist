---
layout: default
title: EKS Checklist Docs
---

# EKS Checklist Docs

이 문서는 EKS 클러스터 점검을 위한 런북을 제공합니다.

📌 **런북 목록**

{% assign categories = site.pages | map: "dir" | uniq | sort %}
{% for category in categories %}
  {% if category != "/" and category != "" %}
  ### {{ category | replace: "/", "" | capitalize }}
  {% assign pages = site.pages | where: "dir", category %}
  {% for page in pages %}
  - [{{ page.title }}]({{ page.url | relative_url }})
  {% endfor %}
  {% endif %}
{% endfor %}
