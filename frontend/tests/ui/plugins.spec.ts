import type { PluginDescriptor, PluginMarketplacePlugin } from "../../features/admin/core/types";
import { test, expect, capture } from "./harness";
import type { MockAPI } from "./network";

const marketplaceEntry: PluginMarketplacePlugin = {
  installed: false, update_available: false,
  plugin: {
    id: "tokenhub.openai-codex", name: "Codex Marketplace Provider", version: "1.0.0", source: "marketplace",
    kinds: [], placements: [], capabilities: [],
    marketplace: { summary: "Synthetic marketplace provider for UI acceptance.", categories: ["provider", "subscription"], publisher: { id: "ui-fixture", name: "UI Fixture Publisher", verified: true } },
  },
};

function plugins(api: MockAPI, entries: PluginMarketplacePlugin[], website = "", installed: () => PluginDescriptor[] = () => []) {
  api.define("GET", "/api/admin/plugins", () => ({ json: { data: structuredClone(installed()) } }));
  api.respond("GET", "/api/admin/plugin-marketplace", { data: { available: true, plugins: entries } });
  api.respond("GET", "/api/admin/provider-adapters", { data: [] });
  api.respond("GET", "/api/admin/plugin-chain", { data: { hooks: [] } });
  api.respond("GET", "/api/admin/plugin-ui-manifest", { data: [] });
  api.respond("GET", "/api/admin/plugin-actions", { data: [] });
  api.respond("GET", "/api/admin/plugin-background-jobs", { data: [], runs: [] });
  api.respond("GET", "/api/admin/resources/settings", { data: [{ id: "cfg_gateway", name: "Gateway", fields: { plugin_marketplace_url: website } }] });
}

for (const mobile of [false, true]) {
  test(`plugins marketplace-only-details${mobile ? "-mobile" : ""}`, async ({ page, api }, testInfo) => {
    if (mobile) await page.setViewportSize({ width: 390, height: 844 });
    plugins(api, [marketplaceEntry]);
    await page.goto("/plugins");
    await page.getByRole("tab", { name: "浏览插件", exact: true }).click();
    const row = page.locator(".plugin-browse-row").filter({ hasText: "Codex Marketplace Provider" });
    await expect(row.getByText("Provider 集成", { exact: true })).toBeVisible();
    await expect(page.getByRole("link", { name: "浏览插件市场", exact: true })).toHaveCount(0);
    await capture(page, testInfo, row, `plugins-browse${mobile ? "-mobile" : ""}`, "市场插件分类和操作");
    await row.getByRole("button", { name: "详情", exact: true }).click();
    await expect(page).toHaveURL(/\/plugins\/tokenhub\.openai-codex$/);
    await expect(page.getByRole("heading", { name: "Codex Marketplace Provider", exact: true })).toBeVisible();
    await expect(page.getByText("Synthetic marketplace provider for UI acceptance.", { exact: true }).first()).toBeVisible();
    await expect(page.locator(".plugin-overview-facts").getByText("未安装", { exact: true })).toBeVisible();
    expect(api.calls.some(call => call.path.endsWith("/detail"))).toBe(false);
    await capture(page, testInfo, page.locator(".plugin-detail-view"), `plugins-marketplace-detail${mobile ? "-mobile" : ""}`, "未安装插件详情", "viewport");
  });
}

test("plugins empty-marketplace-with-configured-website", async ({ page, api }, testInfo) => {
  plugins(api, [], "https://marketplace.example.test");
  await page.goto("/plugins");
  await page.getByRole("tab", { name: "浏览插件", exact: true }).click();
  await expect(page.getByText("暂无可安装插件", { exact: true })).toBeVisible();
  await expect(page.getByRole("link", { name: "浏览插件市场", exact: true })).toHaveAttribute("href", "https://marketplace.example.test/");
  await capture(page, testInfo, page.locator(".plugin-manager-topbar"), "plugins-empty-marketplace", "空市场和已配置网站", "viewport");
});

for (const action of ["install", "update", "uninstall", "rollback"] as const) {
  test(`plugins package-${action}-refresh`, async ({ page, api }, testInfo) => {
    const descriptor: PluginDescriptor = {
      id: "tokenhub.ui-package", name: "UI Lifecycle Package", version: "2.0.0", source: "local_file",
      status: "disabled", installed: true, kinds: ["extension"], placements: ["gateway_chain"], capabilities: [],
      rollback_available: true, rollback_version: "1.0.0",
      distribution: { download_url: "https://plugins.example.test/package.zip", checksum_sha256: "a".repeat(64) },
    };
    let installed = action === "install" ? [] : [descriptor];
    const entries: PluginMarketplacePlugin[] = [{ plugin: { ...descriptor, installed: false }, installed: action !== "install", update_available: false }];
    plugins(api, entries, "", () => installed);
    const endpoint = action === "install" ? "/api/admin/plugins/install" : action === "uninstall" ? `/api/admin/plugin-packages/${descriptor.id}` : `/api/admin/plugins/${descriptor.id}/${action}`;
    api.define(action === "uninstall" ? "DELETE" : "POST", endpoint, input => {
      if (action === "install" || action === "update") expect(input.body).toMatchObject({ download_url: descriptor.distribution!.download_url, checksum_sha256: "a".repeat(64) });
      else expect(input.body).toBeUndefined();
      const version = action === "rollback" ? "1.0.0" : "3.0.0";
      installed = action === "uninstall" ? [] : [{ ...descriptor, version, rollback_available: false }];
      entries[0].installed = installed.length > 0;
      return { json: { data: { plugin: { ...descriptor, version }, plugin_id: descriptor.id, rollback_version: version, restart_required: false } } };
    });
    await page.goto("/plugins");
    if (action === "install") {
      await page.getByRole("button", { name: "安装本地插件", exact: true }).click();
      await page.getByLabel("下载 URL", { exact: true }).fill(descriptor.distribution!.download_url!);
      await page.getByLabel("SHA-256 校验", { exact: true }).fill("a".repeat(64));
      await page.getByRole("button", { name: "安装插件", exact: true }).click();
    } else {
      const labels = { update: "更新", uninstall: "卸载", rollback: "回滚" };
      await page.getByRole("button", { name: labels[action], exact: true }).click();
    }
    await expect.poll(() => api.calls.filter(call => call.method === "GET" && call.path === "/api/admin/plugins").length).toBe(2);
    await page.getByRole("tab", { name: "已安装插件", exact: true }).click();
    const row = page.locator(".plugin-installed-row").filter({ hasText: descriptor.name });
    if (action === "uninstall") await expect(row).toHaveCount(0);
    else await expect(row).toContainText(action === "rollback" ? "1.0.0" : "3.0.0");
    await capture(page, testInfo, page.locator(".plugins-view"), `package-${action}-installed`, "插件操作后安装列表已刷新", "viewport");
    await page.getByRole("tab", { name: "浏览插件", exact: true }).click();
    const available = page.locator(".plugin-browse-row").filter({ hasText: descriptor.name });
    await expect(available).toHaveCount(action === "uninstall" ? 1 : 0);
  });
}
