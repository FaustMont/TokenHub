# Text embeddings

TokenHub exposes `POST /v1/embeddings` with the normal API-key permissions, routing, limits and request audit. Send `model` and either one text or an array of independent texts. A batch returns one dense vector per input, preserving `index`; malformed upstream counts, indices and dimensions are errors.

Optional parameters are `dimensions`, `encoding_format` (`float` or `base64`), `input_type` (`query` or `document`), `task`, `normalized`, `truncation`, `late_chunking` and `user`. Availability depends on the selected protocol. Unsupported options are rejected rather than silently discarded. Sparse, quantized and multimodal inputs/outputs and asynchronous Batch jobs are outside this endpoint's current scope. Token ID inputs remain available on compatible upstreams.

## Provider configuration

In the Provider's advanced settings, open **Text embedding settings**. `embedding_protocol` selects `openai`, `cohere`, `jina`, `voyage`, `dashscope` or `tei`; Gemini uses its native adapter. An empty value uses the catalog default or OpenAI compatibility. Custom providers must select their actual protocol. The path is appended to the Provider Base URL; `embedding_path` overrides that relative endpoint without changing the host or credentials.

| Protocol | Example Base URL ending | Default path | Notes |
| --- | --- | --- | --- |
| OpenAI compatibility | `/v1` | `/embeddings` | OpenAI, SiliconFlow, compatible vLLM/Xinference and Alibaba endpoints |
| Cohere | `/v2` | `/embed` | `input_type` or `task` required; dense float output |
| Voyage / Jina | `/v1` | `/embeddings` | Task and dimension fields are mapped to the native names |
| DashScope native | `/api/v1` | `/services/embeddings/text-embedding/text-embedding` | Text input, dense output; gateway limit of 10 inputs per native call |
| TEI native | server root | `/embed` | Text batches; no invented usage when the server omits it |
| Gemini | configured Gemini base | `:embedContent` / `:batchEmbedContents` | Independent texts remain independent requests |

Use the provider's documented regional URL. A catalog listing or passing HTTP connection test does not establish that a model supports every parameter. Native protocol fixtures validate conversion locally; real account, regional and model availability require upstream verification.

## Vector compatibility

By default, failover stays within the same Provider and upstream model. To permit verified equivalent deployments, configure `embedding_spaces` as a JSON mapping from **upstream model ID** to an administrator-confirmed space ID, for example `{"my-embedding-model":"space-v1"}`. Every participating deployment must use equivalent model versions and encoding behavior. Equal dimensions alone are insufficient. Changing spaces requires rebuilding the application's stored vector index; TokenHub never edits that external index.

Missing upstream usage remains unreported in metering evidence, distinct from reported zero. Estimates are not substituted for measured tokens. Provider costs and tenant prices remain separate.

```json
{"model":"public-embedding","input":["first document","second document"],"dimensions":1024,"encoding_format":"float"}
```

OpenAI SDK clients can use `client.embeddings.create`. Configure the public model name and TokenHub base URL in Dify or the relevant LangChain/LlamaIndex integration. Actual client/version acceptance is separate from local protocol tests.

Model discovery prefers explicit `type`, `modality` or `model_type`. When absent, established BGE, GTE, E5, Voyage and sentence-transformer names are recognized as embedding candidates; reranker names take precedence. Inferred type still does not prove deployed capability.
