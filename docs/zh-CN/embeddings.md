# 文本 Embedding

TokenHub 通过 `POST /v1/embeddings` 提供文本向量，复用 API Key 权限、路由、限流和请求审计。传入 `model` 与单条文本或独立文本数组；批量输入逐条返回向量并保留 `index`。上游返回数量、索引或维度不符时报告错误。

可选参数：`dimensions`、`encoding_format`（`float` / `base64`）、`input_type`（`query` / `document`）、`task`、`normalized`、`truncation`、`late_chunking`、`user`。具体支持取决于上游协议；不支持的参数明确拒绝，不静默丢弃。本接口暂不支持稀疏、量化、多模态和异步 Batch；兼容上游仍可接收 token ID 输入。

## 配置上游

在 Provider 高级设置中打开“文本 Embedding 配置”。协议可选 `openai`、`cohere`、`jina`、`voyage`、`dashscope`、`tei`；Gemini 走原生适配器。留空采用目录默认值或 OpenAI 兼容协议，自定义渠道应明确选择实际协议。

| 协议 | Base URL 结尾示例 | 默认路径 | 说明 |
| --- | --- | --- | --- |
| OpenAI 兼容 | `/v1` | `/embeddings` | OpenAI、SiliconFlow、兼容的 vLLM/Xinference、阿里兼容接口 |
| Cohere | `/v2` | `/embed` | 必须指定 input_type 或 task，返回稠密向量 |
| Voyage / Jina | `/v1` | `/embeddings` | 转换任务和维度参数 |
| DashScope 原生 | `/api/v1` | `/services/embeddings/text-embedding/text-embedding` | 稠密文本向量，网关每次限制 10 条 |
| TEI 原生 | 服务根地址 | `/embed` | 不伪造上游未报告的用量 |
| Gemini | Gemini 基础地址 | `:embedContent` / `:batchEmbedContents` | 独立文本保持独立，不合并批次 |

`embedding_path` 可覆盖追加到 Base URL 的路径，不改变主机或凭据。地域地址以供应商文档为准。目录收录或连接测试通过不代表模型支持全部参数。本地协议测试与真实账号、地域、型号验证分别记录。

## 向量空间

默认只在相同 Provider、相同上游模型之间故障切换。确认部署兼容后，可在 `embedding_spaces` 填写上游模型到空间标识的 JSON，例如 `{"my-embedding-model":"space-v1"}`。相同标识意味着管理员确认模型版本与编码行为兼容；仅维度相同不够。更换空间需要业务应用重建已有向量索引，TokenHub 不修改外部索引。

计量证据区分未报告与明确零用量，不将估算冒充实测。上游成本和租户收费分别配置。

```json
{"model":"public-embedding","input":["第一条文本","第二条文本"],"dimensions":1024,"encoding_format":"float"}
```

OpenAI SDK 使用 `client.embeddings.create`。Dify、LangChain/LlamaIndex 使用对应集成并配置 TokenHub 地址和对外模型名；客户端实际版本验收独立于本地协议测试。

模型发现优先采用上游明确声明的 type、modality 或 model_type；缺失时识别常见 BGE、GTE、E5、Voyage 和 sentence-transformer 名称，重排名称优先判为 rerank。名称推断仍不代表部署能力已经验证。

缓存查询和路由前会检查全部启用的路由定义、资源覆盖及可能的 Provider 回退来源；暂时不健康或冷却中的节点也纳入空间判断。健康变化和加权排序不能切换空间。配置冲突返回 `409 embedding_space_conflict`；请移除不兼容线路，或在确认兼容后填写相同的空间标识。
