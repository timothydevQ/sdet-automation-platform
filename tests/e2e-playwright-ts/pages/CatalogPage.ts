import { Page, expect } from "@playwright/test";

export class CatalogPage {
  constructor(private page: Page) {}

  async goto() {
    await this.page.goto("/");
  }

  async addToCart(sku: string) {
    const btn = this.page.getByTestId(`add-${sku}`);
    await btn.click();
    // Wait for the success message to confirm the API call completed
    await expect(this.page.getByTestId("catalog-msg")).toBeVisible({ timeout: 10000 });
  }

  async waitForProducts() {
    await this.page.getByTestId(/^product-/).first().waitFor({ timeout: 15000 });
  }
}
