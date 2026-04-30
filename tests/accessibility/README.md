# Accessibility tests

axe-core checks live in the Playwright suite at
`tests/e2e-playwright-ts/specs/accessibility.spec.ts`.

Run them in isolation:
```
cd tests/e2e-playwright-ts && npx playwright test specs/accessibility.spec.ts
```

Standards covered: WCAG 2 A and AA. Critical violations fail CI; serious are
reported but not blocking yet (see `docs/test-strategy.md`).
