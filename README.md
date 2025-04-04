# lolidump

用于导出 ManyACG 数据的 CLI 工具.

WIP

## 配置

```toml
[database]
database = "manyacg"
host = "127.0.0.1"
port = 27017
user = "krau"
password = "qwqowo"

[dest]
type = "meilisearch"

[dest.meilisearch]
host = "http://127.0.0.1:7700"
key = ""
index = "manyacg"
[dest.meilisearch.embedder]
name = "default"
source = "ollama"
model = "bge-m3"
document_template = ""
url = "http://127.0.0.1:11434/api/embed"
dimensions = 1024
