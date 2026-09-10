import { afterEach, describe, expect, it } from "vitest";

import {
  type AppLanguage,
  countRatioWithUnit,
  countWithLabel,
  countWithUnit,
  millisecondsText,
  providerSaveMessage,
  routeAttemptCountText,
  setActiveLanguage,
  translateGeneratedText,
} from "./runtime";
import { gatewayLanguageLabel } from "../views/gateway-docs-ui";

describe("localized count ratios", () => {
  afterEach(() => setActiveLanguage("en"));

  it.each<[AppLanguage, string]>([
    ["zh-CN", "2/3 个可引入模型"],
    ["en", "2/3 importable models"],
    ["ja", "2/3 件の取り込み可能モデル"],
    ["ru", "2/3 модели для импорта"],
  ])("preserves the current and total counts in %s", (language, expected) => {
    setActiveLanguage(language);

    expect(countRatioWithUnit(2, 3, "个可引入模型", "importable model", "件の取り込み可能モデル", "importable models")).toBe(expected);
  });

  it("formats Russian countWithLabel, countWithUnit, and routeAttemptCountText correctly", () => {
    setActiveLanguage("ru");
    expect(countWithLabel(1, "个模型")).toBe("1 модель");
    expect(countWithLabel(2, "个模型")).toBe("2 модели");
    expect(countWithLabel(5, "个模型")).toBe("5 моделей");
    expect(countWithLabel(11, "个模型")).toBe("11 моделей");
    expect(countWithLabel(12, "个模型")).toBe("12 моделей");
    expect(countWithLabel(21, "个模型")).toBe("21 модель");
    expect(countWithLabel(22, "个模型")).toBe("22 модели");
    expect(countWithLabel(25, "个模型")).toBe("25 моделей");

    expect(countWithLabel(1, "个 Key")).toBe("1 ключ");
    expect(countWithLabel(2, "个 Key")).toBe("2 ключа");
    expect(countWithLabel(5, "个 Key")).toBe("5 ключей");
    expect(countWithLabel(21, "个 Key")).toBe("21 ключ");

    expect(countWithLabel(1, "个项目")).toBe("1 проект");
    expect(countWithLabel(2, "个项目")).toBe("2 проекта");
    expect(countWithLabel(5, "个项目")).toBe("5 проектов");

    expect(countWithUnit(1, "人", "member", "人")).toBe("1 участник");
    expect(countWithUnit(2, "人", "member", "人")).toBe("2 участника");
    expect(countWithUnit(5, "人", "member", "人")).toBe("5 участников");
    expect(countWithUnit(1, "человек", "person", "人")).toBe("1 человек");
    expect(countWithUnit(2, "человек", "person", "人")).toBe("2 человека");
    expect(countWithUnit(5, "человек", "person", "人")).toBe("5 человек");

    expect(countWithUnit(1, "个用户", "user", "ユーザー")).toBe("1 пользователь");
    expect(countWithUnit(2, "个用户", "user", "ユーザー")).toBe("2 пользователя");
    expect(countWithUnit(5, "个用户", "user", "ユーザー")).toBe("5 пользователей");

    expect(countWithUnit(1, "个OpenAI上游模型", "OpenAI upstream model", "OpenAI 上流モデル")).toBe("1 OpenAI вышестоящая модель");
    expect(countWithUnit(2, "个OpenAI上游模型", "OpenAI upstream model", "OpenAI 上流モデル")).toBe("2 OpenAI вышестоящие модели");
    expect(countWithUnit(5, "个OpenAI上游模型", "OpenAI upstream model", "OpenAI 上流モデル")).toBe("5 OpenAI вышестоящих моделей");

    expect(countWithUnit(1, "条启用线路", "active route", "件の有効ルート")).toBe("1 активный маршрут");
    expect(countWithUnit(2, "条启用线路", "active route", "件の有効ルート")).toBe("2 активных маршрута");
    expect(countWithUnit(5, "条启用线路", "active route", "件の有効ルート")).toBe("5 активных маршрутов");

    expect(routeAttemptCountText(-1)).toBe("-1 попытка");
    expect(routeAttemptCountText(0)).toBe("0 попыток");
    expect(routeAttemptCountText(1)).toBe("1 попытка");
    expect(routeAttemptCountText(2)).toBe("2 попытки, с fallback");
    expect(routeAttemptCountText(3)).toBe("3 попытки, с fallback");
    expect(routeAttemptCountText(5)).toBe("5 попыток, с fallback");
    expect(routeAttemptCountText(21)).toBe("21 попытка, с fallback");
    expect(routeAttemptCountText(22)).toBe("22 попытки, с fallback");
    expect(routeAttemptCountText(25)).toBe("25 попыток, с fallback");

    expect(countRatioWithUnit(0, 3, "个可引入模型", "importable model", "件の取り込み可能モデル", "importable models")).toBe("0/3 моделей для импорта");
    expect(countRatioWithUnit(1, 3, "个可引入模型", "importable model", "件の取り込み可能モデル", "importable models")).toBe("1/3 модель для импорта");
    expect(countRatioWithUnit(2, 3, "个可引入模型", "importable model", "件の取り込み可能モデル", "importable models")).toBe("2/3 модели для импорта");
    expect(countRatioWithUnit(5, 5, "个可引入模型", "importable model", "件の取り込み可能モデル", "importable models")).toBe("5/5 моделей для импорта");
    expect(countRatioWithUnit(2, 3, "custom_item", "custom_item", "custom_item")).toBe("2/3 custom_items");

    expect(millisecondsText(150)).toBe("150 мс");
    expect(translateGeneratedText("已提交审批：API Key 发放", "ru")).toBe("Заявка на согласование отправлена: API Key 发放");

    expect(providerSaveMessage(true, true, 2, "OpenAI")).toBe("Provider обновлен, ресурс аккаунта создан, импортировано: 2 OpenAI вышестоящие модели");
    expect(providerSaveMessage(false, false, 1, "OpenAI")).toBe("Provider создан, импортировано: 1 OpenAI вышестоящая модель");
    expect(providerSaveMessage(false, false, 0, "OpenAI")).toBe("Provider создан");

    expect(gatewayLanguageLabel("zh-CN")).toBe("中文");
    expect(gatewayLanguageLabel("ja")).toBe("日本語");
    expect(gatewayLanguageLabel("ru")).toBe("Русский");
    expect(gatewayLanguageLabel("en")).toBe("English");
  });
});
