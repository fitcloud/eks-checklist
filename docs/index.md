# 기능별 런북

{% for category in site.collections %}
## {{ category.label | capitalize }}

{% for page in category.docs %}
- [{{ page.title | default: page.name | replace: ".md", "" }}]({{ page.url | relative_url }})
{% endfor %}

{% endfor %}
