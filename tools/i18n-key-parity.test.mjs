// Key parity between English, Japanese, and Russian admin console dictionaries.

import assert from "node:assert/strict";
import { readFileSync } from "node:fs";
import { dirname, join } from "node:path";
import { fileURLToPath } from "node:url";
import { describe, it } from "node:test";

const i18nDir = join(
  dirname(fileURLToPath(import.meta.url)),
  "..",
  "frontend",
  "features",
  "admin",
  "i18n",
);

const TYPE_ANNOTATIONS = [
  ': Record<"en" | "ja" | "ru", Record<string, string>>',
  ': Record<"en" | "ja", Record<string, string>>',
  ": Record<string, string>",
  ' satisfies Record<"en" | "ja" | "ru", Record<string, string>>',
  ' satisfies Record<"en" | "ja", Record<string, string>>',
  " as const",
];

export function stripTypeAnnotations(source) {
  let javascript = source;
  for (const annotation of TYPE_ANNOTATIONS) javascript = javascript.replaceAll(annotation, "");
  javascript = javascript.replace(/:\s*Record<[^;{]+>/g, "");
  javascript = javascript.replace(/\s*satisfies\s+Record<[^;]+>(?=\s*;?)/g, "");
  return javascript;
}

async function loadDictionarySource(file) {
  const source = readFileSync(join(i18nDir, file), "utf8");
  assert.ok(
    !/^\s*import\s/m.test(source),
    `${file} now has an import; this loader only handles self-contained dictionaries.`,
  );
  const javascript = stripTypeAnnotations(source);
  assert.ok(
    !javascript.includes("Record<"),
    `${file} uses a type annotation this loader does not strip.`,
  );
  return import(`data:text/javascript;base64,${Buffer.from(javascript).toString("base64")}`);
}

const { enTranslations } = await loadDictionarySource("en.tsx");
const { jaTranslations } = await loadDictionarySource("ja.tsx");
const { ruTranslations } = await loadDictionarySource("ru.tsx");
const { adminWorkflowTranslations } = await loadDictionarySource("admin-workflows.tsx");
const { apiKeyUsageTranslations } = await loadDictionarySource("api-key-usage.tsx");
const { auditFilterTranslations } = await loadDictionarySource("audit-filters.tsx");
const { dbEvolutionTranslations } = await loadDictionarySource("db-evolution.tsx");
const { loginHomeTranslations } = await loadDictionarySource("login-home.tsx");
const { modelGovernanceTranslations } = await loadDictionarySource("model-governance.tsx");
const { gatewayDocsTranslations } = await loadDictionarySource("gateway-docs.tsx");
const { playgroundTranslations } = await loadDictionarySource("playground.tsx");
const { providerConnectionTranslations } = await loadDictionarySource("provider-connection.tsx");
const { providerMonitoringTranslations } = await loadDictionarySource("provider-monitoring.tsx");
const { routingTranslations } = await loadDictionarySource("routing.tsx");
const { codexImageTranslations } = await loadDictionarySource("codex-image.tsx");
const { scopedRoutingPolicyTranslations } = await loadDictionarySource("scoped-routing-policy.tsx");
const { securityTranslations } = await loadDictionarySource("security.tsx");
const { usageTranslations } = await loadDictionarySource("usage.tsx");
const { notificationTranslations } = await loadDictionarySource("notifications.tsx");
const syntheticModule = await loadDictionarySource("synthetic-dns.tsx");
const syntheticDNSTranslations = syntheticModule.default ?? syntheticModule.syntheticDNSTranslations;

const submodules = [
  ["admin-workflows.tsx", adminWorkflowTranslations],
  ["api-key-usage.tsx", apiKeyUsageTranslations],
  ["audit-filters.tsx", auditFilterTranslations],
  ["db-evolution.tsx", dbEvolutionTranslations],
  ["routing.tsx", routingTranslations],
  ["codex-image.tsx", codexImageTranslations],
  ["scoped-routing-policy.tsx", scopedRoutingPolicyTranslations],
  ["model-governance.tsx", modelGovernanceTranslations],
  ["gateway-docs.tsx", gatewayDocsTranslations],
  ["login-home.tsx", loginHomeTranslations],
  ["provider-connection.tsx", providerConnectionTranslations],
  ["provider-monitoring.tsx", providerMonitoringTranslations],
  ["usage.tsx", usageTranslations],
  ["playground.tsx", playgroundTranslations],
  ["security.tsx", securityTranslations],
  ["notifications.tsx", notificationTranslations],
  ["synthetic-dns.tsx", syntheticDNSTranslations],
];

// Mirrors the merge order in translations.tsx.
const merged = {
  en: {
    ...enTranslations,
    ...adminWorkflowTranslations.en,
    ...apiKeyUsageTranslations.en,
    ...auditFilterTranslations.en,
    ...dbEvolutionTranslations.en,
    ...routingTranslations.en,
    ...codexImageTranslations.en,
    ...scopedRoutingPolicyTranslations.en,
    ...modelGovernanceTranslations.en,
    ...gatewayDocsTranslations.en,
    ...loginHomeTranslations.en,
    ...providerConnectionTranslations.en,
    ...providerMonitoringTranslations.en,
    ...usageTranslations.en,
    ...playgroundTranslations.en,
    ...securityTranslations.en,
    ...notificationTranslations.en,
    ...syntheticDNSTranslations.en,
  },
  ja: {
    ...jaTranslations,
    ...adminWorkflowTranslations.ja,
    ...apiKeyUsageTranslations.ja,
    ...auditFilterTranslations.ja,
    ...dbEvolutionTranslations.ja,
    ...routingTranslations.ja,
    ...codexImageTranslations.ja,
    ...scopedRoutingPolicyTranslations.ja,
    ...modelGovernanceTranslations.ja,
    ...gatewayDocsTranslations.ja,
    ...loginHomeTranslations.ja,
    ...providerConnectionTranslations.ja,
    ...providerMonitoringTranslations.ja,
    ...usageTranslations.ja,
    ...playgroundTranslations.ja,
    ...securityTranslations.ja,
    ...notificationTranslations.ja,
    ...syntheticDNSTranslations.ja,
  },
  ru: {
    ...ruTranslations,
    ...adminWorkflowTranslations.ru,
    ...apiKeyUsageTranslations.ru,
    ...auditFilterTranslations.ru,
    ...dbEvolutionTranslations.ru,
    ...routingTranslations.ru,
    ...codexImageTranslations.ru,
    ...scopedRoutingPolicyTranslations.ru,
    ...modelGovernanceTranslations.ru,
    ...gatewayDocsTranslations.ru,
    ...loginHomeTranslations.ru,
    ...providerConnectionTranslations.ru,
    ...providerMonitoringTranslations.ru,
    ...usageTranslations.ru,
    ...playgroundTranslations.ru,
    ...securityTranslations.ru,
    ...notificationTranslations.ru,
    ...syntheticDNSTranslations.ru,
  },
};

function keysMissingFrom(source, target) {
  return Object.keys(source).filter((key) => !Object.hasOwn(target, key));
}

function assertSameKeys(label, en, ja) {
  assert.deepEqual(keysMissingFrom(en, ja), [], `${label}: defined in en, missing from ja`);
  assert.deepEqual(keysMissingFrom(ja, en), [], `${label}: defined in ja, missing from en`);
}

function assertSameThreeWayKeys(label, en, ja, ru) {
  assert.deepEqual(keysMissingFrom(en, ja), [], `${label}: defined in en, missing from ja`);
  assert.deepEqual(keysMissingFrom(ja, en), [], `${label}: defined in ja, missing from en`);
  assert.deepEqual(keysMissingFrom(en, ru), [], `${label}: defined in en, missing from ru`);
  assert.deepEqual(keysMissingFrom(ru, en), [], `${label}: defined in ru, missing from en`);
}

describe("dictionary loading", () => {
  it("parses every dictionary into a non-trivial object", () => {
    const sources = {
      "en.tsx": enTranslations,
      "ja.tsx": jaTranslations,
      "ru.tsx": ruTranslations,
    };
    for (const [name, mod] of submodules) {
      sources[`${name} en`] = mod.en;
      sources[`${name} ja`] = mod.ja;
      sources[`${name} ru`] = mod.ru;
    }
    for (const [label, dictionary] of Object.entries(sources)) {
      const count = Object.keys(dictionary).length;
      assert.ok(count > 0, `${label} parsed to 0 keys`);
      for (const [key, value] of Object.entries(dictionary)) {
        assert.equal(typeof value, "string", `${label}: ${key} is not a string`);
        assert.notEqual(value.trim(), "", `${label}: ${key} has an empty translation`);
      }
    }
  });

  it("strips satisfies and colon Record type annotations without syntax errors", () => {
    const raw = `export const sample = { en: { a: "b" } } satisfies Record< "en" | "ja" | "ru", Record<string, string> >;`;
    const stripped = stripTypeAnnotations(raw);
    assert.ok(!stripped.includes("Record<"));
    assert.equal(stripped.trim(), `export const sample = { en: { a: "b" } };`);
  });
});

describe("single ownership", () => {
  // Scope, stated honestly: this compares the evaluated objects, so it catches a key
  // defined in two different files. A key repeated twice inside one file collapses
  // during evaluation and is invisible here.
  const sourcesByLanguage = {
    en: [
      ["en.tsx", enTranslations],
      ["routing.tsx", routingTranslations.en],
      ["model-governance.tsx", modelGovernanceTranslations.en],
    ],
    ja: [
      ["ja.tsx", jaTranslations],
      ["routing.tsx", routingTranslations.ja],
      ["model-governance.tsx", modelGovernanceTranslations.ja],
    ],
    ru: [
      ["ru.tsx", ruTranslations],
      ["routing.tsx", routingTranslations.ru],
      ["model-governance.tsx", modelGovernanceTranslations.ru],
    ],
  };

  for (const [language, sources] of Object.entries(sourcesByLanguage)) {
    it(`does not define a ${language} key in multiple sources`, () => {
      // A key defined twice silently resolves to whichever source translations.tsx
      // merges last, leaving the other definition as dead text that still reads like it
      // is in use. "示例" sat that way in en.tsx, shadowed by routing.tsx and disagreeing
      // with it ("Examples" against "Example"), which is the trap this catches.
      const owners = new Map();
      const shadowed = [];
      for (const [file, dictionary] of sources) {
        for (const key of Object.keys(dictionary)) {
          const owner = owners.get(key);
          if (owner) shadowed.push(`${key} (${owner} shadowed by ${file})`);
          else owners.set(key, file);
        }
      }
      assert.deepEqual(shadowed, [], `${language} keys defined in more than one source`);
    });
  }
});

describe("key parity across all submodules", () => {
  it("keeps primary dictionaries in step", () => {
    assertSameThreeWayKeys("primary dictionaries (en vs ja vs ru)", enTranslations, jaTranslations, ruTranslations);
  });

  for (const [name, mod] of submodules) {
    it(`keeps ${name} in step across en, ja, and ru`, () => {
      assertSameThreeWayKeys(`${name} (en vs ja vs ru)`, mod.en, mod.ja, mod.ru);
    });
  }

  it("keeps the merged dictionary tx() reads in step", () => {
    assertSameThreeWayKeys("merged translations (en vs ja vs ru)", merged.en, merged.ja, merged.ru);
  });
});
