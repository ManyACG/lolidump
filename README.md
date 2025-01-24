# lolidump

用于导出 ManyACG 数据的 CLI 工具.

WIP

## 配置

```toml
[database]
database = "manyacg"
host = "192.168.31.5"
port = 27017
user = "krau"
password = "qwqowo"

[dest]
type = "meilisearch"

[dest.meilisearch]
host = "http://localhost:7700"
key = "114514"
```