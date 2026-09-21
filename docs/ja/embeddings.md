# テキスト埋め込み

`POST /v1/embeddings` は API キー権限、ルーティング、制限、監査を適用します。`model` と単一テキストまたは独立したテキスト配列を送信します。入力ごとに密ベクトルを返し、`index` を保持します。上流の件数、索引、次元が不正な場合はエラーです。

任意パラメータは `dimensions`、`encoding_format`（float/base64）、`input_type`（query/document）、`task`、`normalized`、`truncation`、`late_chunking`、`user` です。対応はプロトコルによって異なり、未対応の指定は拒否します。疎ベクトル、量子化、マルチモーダル、非同期 Batch は対象外です。互換上流では token ID 入力を利用できます。

## プロバイダー設定

詳細設定のテキスト埋め込み設定で `embedding_protocol` を選択します。`openai`、`cohere`、`jina`、`voyage`、`dashscope`、`tei` に対応し、Gemini は専用アダプターを使います。空欄はカタログ既定値または OpenAI 互換です。カスタム接続では実際のプロトコルを選択してください。

| プロトコル | Base URL の末尾例 | 既定パス | 注意点 |
| --- | --- | --- | --- |
| OpenAI 互換 | `/v1` | `/embeddings` | SiliconFlow、対応 vLLM/Xinference、Alibaba 互換 API など |
| Cohere | `/v2` | `/embed` | input_type または task が必要 |
| Voyage / Jina | `/v1` | `/embeddings` | タスクと次元のパラメータを変換 |
| DashScope | `/api/v1` | `/services/embeddings/text-embedding/text-embedding` | 密ベクトル、1 リクエスト最大 10 入力 |
| TEI | サーバールート | `/embed` | 未報告の使用量を生成しない |
| Gemini | Gemini 基本 URL | `:embedContent` / `:batchEmbedContents` | 独立テキストを結合しない |

`embedding_path` は Base URL に追加するパスを変更し、ホストや認証情報は変更しません。リージョン URL は公式資料で確認してください。カタログや接続テストだけでは全パラメータ対応を保証できません。ローカルプロトコルテストと実アカウント・モデル検証は区別します。

## ベクトル空間

既定のフェイルオーバーは同じ Provider と上流モデルに限定します。互換性確認済みの配備には `embedding_spaces` にモデル ID と空間 ID の JSON 対応表（例 `{"my-embedding-model":"space-v1"}`）を設定できます。モデルのバージョンと符号化動作の互換性を確認してください。同じ次元数だけでは不十分です。空間変更時はアプリケーションで索引を再構築します。TokenHub は外部索引を変更しません。

計量証拠では未報告とゼロを区別し、推定を実測に置き換えません。上流コストとテナント料金は別設定です。

```json
{"model":"public-embedding","input":["最初の文書","次の文書"],"dimensions":1024,"encoding_format":"float"}
```

OpenAI SDK は `client.embeddings.create` を使用します。Dify、LangChain/LlamaIndex では対応する統合に公開モデル名と TokenHub URL を設定してください。実際のクライアントバージョンの受け入れ検証は別途必要です。

モデル検出は明示された type、modality、model_type を優先します。未指定の場合は BGE、GTE、E5、Voyage、sentence-transformer の既知の名前を認識し、再ランキング名を優先して分類します。名前からの推定は配備能力の検証ではありません。

キャッシュ照会やルーティング前に、有効なルート定義、リソースの上書き、Provider へのフォールバック候補が同じ空間か確認します。不健全・クールダウン中のノードも対象です。健康状態や重み付き順序で空間は変わりません。競合は `409 embedding_space_conflict` を返します。非互換ルートを削除するか、互換性確認後に同じ空間 ID を設定してください。

フックはテキストを書き換えられますが、入力数、次元、符号化、タスクの意味を保持する必要があります。キャッシュ・プラグイン・後処理の出力は元の契約で検証します。キャッシュプラグインには空間・要求・呼び出し元に結び付くホスト生成 cache_key を渡します。ヒット時は保存済みキーを cache_key 書き込みで返し、欠落・不一致はミスとして扱います。書き込み側にも同じキーを渡すため、既存プラグインは対応が必要です。
