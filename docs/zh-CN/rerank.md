# 文本重排

`POST /v1/rerank` 根据查询对候选文档排序，接入 API Key、模型权限、限流、适用安全策略、插件、路由和审计。对外模型必须为 `rerank`，渠道协议及上下游价格配置完整后才能正式调用。

```json
{"model":"public-reranker","query":"如何续费证书？","documents":["证书续费说明","无关内容"],"top_n":1,"return_documents":true}
```

响应包含 `model` 和 `results`，结果保留原始 `index`、`relevance_score`，可选回显 `document.text`。重复文本按索引区分。不修改分数范围，不拆批后假定全局排序等价；空查询/文档、非法 top_n、不支持的参数、上游缺项或错误索引均明确报错。本次仅支持文本，不包含图片、视频和异步任务。

## 配置渠道

在“高级 → 文本重排配置”选择服务实际提供的协议。路径追加到 Base URL，可用 `rerank_path` 显式覆盖。

| 协议 | Base URL 结尾 | 默认路径 | 场景 |
| --- | --- | --- | --- |
| jina | `/v1` | `/rerank` | Jina、SiliconFlow、兼容的 vLLM/Xinference；本地 BGE 必须由服务框架提供 API |
| cohere | `/v2` | `/rerank` | Cohere 文本重排 |
| voyage | `/v1` | `/rerank` | top_n 转换为 top_k |
| qwen | 对应地域兼容 API 地址 | `/reranks` | 阿里 qwen3-rerank |
| dashscope | `/api/v1` | `/services/rerank/text-rerank/text-rerank` | 阿里 gte-rerank-v2 原生封装 |
| tei | 服务根地址 | `/rerank` | TEI，documents 转换为 texts |

不能从模型名称推断部署协议。Xinference 应填写实际 MODEL_UID。阿里官方资料存在不同 URL 前缀写法，地域和模型可用性也不同，必须核实实际端点；不会把 qwen-rerank 自动替换为其他模型。

适用的 jina/qwen/dashscope 协议支持 instruction 并映射原生字段；voyage/tei 支持 truncation。具体模型限制仍由上游校验，不支持的协议参数明确拒绝。

## 价格与既有配置

按 Token 计量的模型使用输入单价；免费检索价格需要显式确认，空配置不视为免费。Cohere 的 search_units 使用上下游各自的“搜索单元价格 USD/次”，API 对应 `metadata.search_unit_price_usd`；明确免费 Token 价格使用 `metadata.retrieval_pricing_confirmed="true"`。保存渠道成本时确认该渠道价格。上游成本与租户收费独立计算，计量证据区分未报告与实测零值。

目录中无法调用的模型保留展示并标记未支持，不可用于新增路由；新发布检查能力和价格。升级不会自动停用既有可工作路由；更换模型/渠道或重新发布时执行新校验。此前误判的 BGE 可显式改为 rerank 或重新发现；不批量修改，也不重写历史账单。

发布前，管理员可用会话调用 `POST /api/admin/playground/rerank`：

```json
{"provider_id":"configured-provider","request":{"model":"实际上游模型或UID","query":"测试查询","documents":["候选文本"]}}
```

可选 resource_id；结果返回 response、usage_evidence、pricing_status。此入口验证管理员的上游配置，不替代租户权限和计费验收；审计保存模型与计量事实，不记录测试查询/文档全文。

## 验证边界

本地测试覆盖代表协议、请求/响应转换、网关鉴权、结果索引和搜索单元费用；配置界面使用浏览器场景验证。生产账号权限、地域端点和具体服务版本仍需真实上游验证。Dify、LangChain、LlamaIndex 使用对应 HTTP/重排集成；OpenAI SDK 没有原生 rerank 方法。目录存在模型不等于客户端和供应商已完成验收。
