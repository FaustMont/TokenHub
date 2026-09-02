import { afterEach, describe, expect, it } from "vitest";

import {
  type AppLanguage,
  countRatioWithUnit,
  countWithLabel,
  millisecondsText,
  routeAttemptCountText,
  setActiveLanguage,
  translateGeneratedText,
} from "./runtime";

describe("localized count ratios", () => {
  afterEach(() => setActiveLanguage("en"));

  it.each<[AppLanguage, string]>([
    ["zh-CN", "2/3 个可引入模型"],
    ["en", "2/3 importable models"],
    ["ja", "2/3 件の取り込み可能モデル"],
    ["ru", "2/3 моделей для импорта"],
  ])("preserves the current and total counts in %s", (language, expected) => {
    setActiveLanguage(language);

    expect(countRatioWithUnit(2, 3, "个可引入模型", "importable model", "件の取り込み可能モデル", "importable models")).toBe(expected);
  });

  it("formats Russian countWithLabel and routeAttemptCountText correctly", () => {
    setActiveLanguage("ru");
    expect(countWithLabel(5, "个模型")).toBe("5 моделей");
    expect(routeAttemptCountText(1)).toBe("1 попытка");
    expect(routeAttemptCountText(3)).toBe("3 попыток, с fallback");
    expect(millisecondsText(150)).toBe("150 мс");
    expect(translateGeneratedText("已提交审批：API Key 发放", "ru")).toBe("Заявка на согласование отправлена: API Key 发放");
  });
});
