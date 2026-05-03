import { Page, expect } from "@playwright/test";

export class CartPage {
  constructor(private page: Page) {}

  async goto() {
    await this.page.goto("/cart");
    // Wait for cart to finish loading
    await this.page.waitForLoadState("networkidle");
  }

  async checkout() {
    await this.page.getByTestId("checkout-btn").click();
  }

  async expectItemCount(n: number) {
    if (n === 0) {
      await expect(this.page.getByTestId(/^cart-item-/)).toHaveCount(0, { timeout: 10000 });
    } else {
      await expect(this.page.getByTestId(/^cart-item-/).first()).toBeVisible({ timeout: 10000 });
      await expect(this.page.getByTestId(/^cart-item-/)).toHaveCount(n, { timeout: 10000 });
    }
  }
}
