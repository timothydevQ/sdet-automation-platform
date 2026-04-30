import { Page, expect } from "@playwright/test";

export class CartPage {
  constructor(private page: Page) {}

  async goto() {
    await this.page.goto("/cart");
  }

  async checkout() {
    await this.page.getByTestId("checkout-btn").click();
  }

  async expectItemCount(n: number) {
    const items = this.page.getByTestId(/^cart-item-/);
    await expect(items).toHaveCount(n);
  }
}
