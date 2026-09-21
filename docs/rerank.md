# Text reranking

`POST /v1/rerank` scores candidate documents for a query. It uses TokenHub API keys, model permissions, request limits, applicable security policies, gateway hooks, routing and audit. A model must be published as `rerank` with a supported Provider protocol and configured tenant/provider prices.

```json
{"model":"public-reranker","query":"How do I renew a certificate?","documents":["Renewal instructions","An unrelated document"],"top_n":1,"return_documents":true}
```

The response contains `model` and `results` with original `index`, `relevance_score` and optional `document.text`. Duplicate texts remain distinct by index. Scores are not rescaled or comparable across models. The gateway does not split a listwise request into batches. Empty queries/documents, invalid `top_n`, unsupported fields, invalid upstream indices and missing results are errors. Only text is supported; images, video and asynchronous jobs are outside this release.

## Configure a provider

In **Advanced → Text rerank settings**, choose the protocol that the deployed server actually implements. The endpoint path is appended to the Base URL and may be explicitly overridden with `rerank_path`.

| Protocol | Base URL ending | Default endpoint | Examples |
| --- | --- | --- | --- |
| `jina` | `/v1` | `/rerank` | Jina, SiliconFlow, compatible vLLM/Xinference deployments; local BGE requires a serving API |
| `cohere` | `/v2` | `/rerank` | Cohere text rerank |
| `voyage` | `/v1` | `/rerank` | Voyage; maps `top_n` to `top_k` |
| `qwen` | regional compatible API base | `/reranks` | Alibaba `qwen3-rerank` |
| `dashscope` | `/api/v1` | `/services/rerank/text-rerank/text-rerank` | Alibaba `gte-rerank-v2` native envelope |
| `tei` | server root | `/rerank` | Hugging Face TEI; maps documents to `texts` |

A bare model name does not identify the serving protocol. Xinference uses the deployed model UID. Alibaba URLs and availability vary by region; its documentation contains differing URL prefixes, so verify the exact deployed endpoint. No automatic alias from `qwen-rerank` to a different model is created.

`instruction` is available on the applicable Jina-compatible/Qwen/DashScope profiles, with native parameter mapping; `truncation` is available on Voyage and TEI. Support ultimately depends on the selected model. Unsupported protocol options fail explicitly.

## Prices, unsupported models and existing configuration

Token-metered models use the configured input price. Explicitly confirm free retrieval prices rather than relying on an empty legacy zero. Cohere's reported `search_units` use the separate **Price per search unit (USD)** on both the upstream inventory and external model. API configuration stores this decimal as `metadata.search_unit_price_usd`; explicit free token pricing uses `metadata.retrieval_pricing_confirmed="true"`. Provider costs and tenant charges are calculated separately and retain metering evidence; missing usage stays unreported rather than becoming measured zero.

Uncallable catalog models remain visible, are marked unsupported and cannot be selected for new routes. New publication validates capability and prices. Existing working routes are not disabled by upgrade; changing their model/provider or republishing applies the new checks. Previously misclassified BGE entries can be explicitly corrected to `rerank` or rediscovered; no bulk changes or historical bill rewrites occur.

An administrator can test before publication using `POST /api/admin/playground/rerank`:

```json
{"provider_id":"configured-provider","request":{"model":"actual-upstream-model-or-uid","query":"example query","documents":["candidate document"]}}
```

Use the administrator session, optionally supply `resource_id`, and inspect `response`, `usage_evidence` and `pricing_status`. This is an administrative upstream test, not a tenant billing/authorization acceptance test; its audit stores model and metering evidence rather than query/documents.

## Validation boundary

Local tests cover representative protocol payloads and responses, gateway authentication, indices and search-unit charges. The configuration screen has a browser fixture. Production account permissions, regional URLs and exact serving versions require real-upstream validation. Dify, LangChain and LlamaIndex require the matching HTTP/rerank integration; OpenAI SDK does not define a native rerank method. Do not infer full client or supplier acceptance from a model appearing in the directory.
