import type { PluginMarketplacePlugin } from "../../features/admin/core/types";
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

function plugins(api: MockAPI, entries: PluginMarketplacePlugin[], website = "") {
  api.respond("GET", "/api/admin/plugins", { data: [] });
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
