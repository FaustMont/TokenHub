import { afterEach, describe, expect, it } from "vitest";

import { setActiveLanguage } from "../i18n/runtime";
import { permissionDeniedMessage, permissionPartialLoadMessage, readAdminError } from "./payloads";

describe("Provider model catalog error localization", () => {
  it.each([
    ["provider_models_authentication_failed", "401 Unauthorized", "上游拒绝了 Provider 凭据，请检查 API Key 或认证配置。"],
    ["provider_models_rate_limited", "429 Too Many Requests", "上游模型目录请求过于频繁，请稍后重试。"],
    ["provider_models_upstream_error", "500 Internal Server Error", "上游模型目录加载失败，请检查 Provider 连接配置后重试。"],
    ["provider_models_request_failed", "Failed to request upstream models", "无法连接上游模型目录，请检查 Provider 地址和网络配置后重试。"],
    ["provider_models_invalid_response", "Upstream models response is invalid", "上游模型目录返回了无法识别的数据，请检查 Provider 兼容性。"],
    ["provider_models_empty", "Upstream did not return any models", "上游模型目录未返回任何模型，请检查 Provider 配置或稍后重试。"],
  ])("localizes %s without exposing the backend message", async (code, backendMessage, expected) => {
    const response = new Response(JSON.stringify({ error: { code, message: backendMessage } }), {
      status: 502,
      headers: { "content-type": "application/json" },
    });

    const message = await readAdminError(response, "Provider 模型加载失败");

    expect(message).toBe(expected);
    expect(message).not.toBe(backendMessage);
  });
});

describe("Permission message localization", () => {
  afterEach(() => setActiveLanguage("zh-CN"));

  it("formats permissionDeniedMessage in Russian and English", () => {
    setActiveLanguage("ru");
    expect(permissionDeniedMessage("провайдерам")).toBe(
      "У этой учетной записи нет разрешения на доступ к провайдерам. Данные вне вашей области доступа скрыты; обратитесь к администратору для изменения роли или участия в проекте.",
    );
    expect(permissionDeniedMessage("")).toBe(
      "У этой учетной записи нет разрешения на доступ к Этот ресурс. Данные вне вашей области доступа скрыты; обратитесь к администратору для изменения роли или участия в проекте.",
    );

    setActiveLanguage("en");
    expect(permissionDeniedMessage("providers")).toBe(
      "This account does not have permission to access providers. Data outside your permission scope is hidden; ask an admin to adjust your role or project membership if needed.",
    );
  });

  it("formats permissionPartialLoadMessage in Russian with comma separator", () => {
    setActiveLanguage("ru");
    expect(permissionPartialLoadMessage(["Ключи", "Маршруты"])).toBe(
      "Скрыто из-за недостаточных прав: Ключи, Маршруты. На этой странице отображается только доступный вам контент.",
    );

    setActiveLanguage("zh-CN");
    expect(permissionPartialLoadMessage(["Key 管理", "路由策略"])).toBe(
      "已隐藏无权限数据：Key 管理、路由策略。当前页面只展示你有权限查看的内容。",
    );
  });
});

