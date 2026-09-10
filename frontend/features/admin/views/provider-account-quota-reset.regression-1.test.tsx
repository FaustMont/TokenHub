import { afterEach, describe, expect, it } from "vitest";

import { setActiveLanguage } from "../i18n/runtime";
import { formatResetCreditExpiry } from "./provider-account-quota-reset";

describe("formatResetCreditExpiry localization", () => {
  afterEach(() => setActiveLanguage("zh-CN"));

  const now = Date.parse("2026-09-10T12:00:00Z");

  it("handles invalid or expired values", () => {
    setActiveLanguage("ru");
    expect(formatResetCreditExpiry(null, now)).toBe("-");
    expect(formatResetCreditExpiry(undefined, now)).toBe("-");
    expect(formatResetCreditExpiry("invalid-date", now)).toBe("-");
    expect(formatResetCreditExpiry(new Date(now - 1000).toISOString(), now)).toBe("Истёк");

    setActiveLanguage("en");
    expect(formatResetCreditExpiry(new Date(now - 1000).toISOString(), now)).toBe("Expired");

    setActiveLanguage("zh-CN");
    expect(formatResetCreditExpiry(new Date(now - 1000).toISOString(), now)).toBe("已过期");
  });

  it("formats day and hour pluralization in Russian", () => {
    setActiveLanguage("ru");
    const plus25h = new Date(now + 25 * 3600 * 1000).toISOString();
    expect(formatResetCreditExpiry(plus25h, now)).toBe("через 1 день 1 час");

    const plus50h = new Date(now + 50 * 3600 * 1000).toISOString();
    expect(formatResetCreditExpiry(plus50h, now)).toBe("через 2 дня 2 часа");

    const plus125h = new Date(now + 125 * 3600 * 1000).toISOString();
    expect(formatResetCreditExpiry(plus125h, now)).toBe("через 5 дней 5 часов");

    const plus21d1h = new Date(now + (21 * 24 + 1) * 3600 * 1000).toISOString();
    expect(formatResetCreditExpiry(plus21d1h, now)).toBe("через 21 день 1 час");
  });

  it("formats hour and minute pluralization in Russian", () => {
    setActiveLanguage("ru");
    const plus61m = new Date(now + 61 * 60 * 1000).toISOString();
    expect(formatResetCreditExpiry(plus61m, now)).toBe("через 1 час 1 минуту");

    const plus122m = new Date(now + 122 * 60 * 1000).toISOString();
    expect(formatResetCreditExpiry(plus122m, now)).toBe("через 2 часа 2 минуты");

    const plus305m = new Date(now + 305 * 60 * 1000).toISOString();
    expect(formatResetCreditExpiry(plus305m, now)).toBe("через 5 часов 5 минут");
  });

  it("formats minutes-only pluralization in Russian", () => {
    setActiveLanguage("ru");
    expect(formatResetCreditExpiry(new Date(now + 1 * 60 * 1000).toISOString(), now)).toBe("через 1 минуту");
    expect(formatResetCreditExpiry(new Date(now + 2 * 60 * 1000).toISOString(), now)).toBe("через 2 минуты");
    expect(formatResetCreditExpiry(new Date(now + 5 * 60 * 1000).toISOString(), now)).toBe("через 5 минут");
    expect(formatResetCreditExpiry(new Date(now + 21 * 60 * 1000).toISOString(), now)).toBe("через 21 минуту");
    expect(formatResetCreditExpiry(new Date(now + 22 * 60 * 1000).toISOString(), now)).toBe("через 22 минуты");
    expect(formatResetCreditExpiry(new Date(now + 25 * 60 * 1000).toISOString(), now)).toBe("через 25 минут");
  });
});
