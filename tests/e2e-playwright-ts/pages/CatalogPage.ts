import { Page } from "@playwright/test";

export class CatalogPage {
  constructor(private page: Page) {}

  async goto() {
    await this.page.goto("/");
  }

  async addToCart(sku: string) {
    await this.page.getByTestId(`add-${sku}`).click();
  }

  async waitForProducts() {
    await this.page.getByTestId(/^product-/).first().waitFor();
  }
}
