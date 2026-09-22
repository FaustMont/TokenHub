import type { Model } from "../core/types";
import { priceMetric } from "./catalog";
import { formatMoney } from "./formatting";

export function modelTokenPriceMetric(model: Model): string {
  const price = model.modality === "embedding" ? model.embedding_price_usd_per_1m : model.input_price_usd_per_1m;
  const retrieval = model.modality === "embedding" || model.modality === "rerank";
  if (retrieval && (price ?? 0) === 0 && model.metadata?.retrieval_pricing_confirmed === "true") {
    return `$${formatMoney(0)}/Mt`;
  }
  return priceMetric(price);
}
