import { render, screen } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { afterEach, describe, expect, it, vi } from "vitest";
import { emptyData } from "../domain/catalog";
import { setActiveLanguage } from "../i18n/runtime";
import { BillingRateCards } from "./billing-rate-cards";

describe("BillingRateCards", () => {
  afterEach(() => {
    setActiveLanguage("en");
    vi.unstubAllGlobals();
  });

  it("localizes an empty billing error response fallback", async () => {
    vi.stubGlobal("fetch", vi.fn(async () => new Response(null, { status: 500 })));
    setActiveLanguage("en");
    render(<BillingRateCards api={{ baseURL: "http://localhost:8080", adminToken: "admin-token" }} data={emptyData()} />);

    await userEvent.click(screen.getByRole("button", { name: "Load published versions" }));

    expect(await screen.findByRole("alert")).toHaveTextContent("Billing operation failed (500)");
  });
});
