import { Page, expect } from "@playwright/test";

export class AdminOrdersPage {
  constructor(private page: Page) {}

  async goto() {
    await this.page.goto("/admin/orders");
  }

  async refund(id: number) {
    await this.page.getByTestId(`refund-${id}`).click();
  }

  async expectAtLeastOneOrder() {
    await expect(this.page.getByTestId(/^order-row-/).first()).toBeVisible();
  }
}
