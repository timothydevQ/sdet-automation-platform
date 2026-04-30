import { Page } from "@playwright/test";

export class CheckoutPage {
  constructor(private page: Page) {}

  async submit({ coupon = "", card = "tok_test_visa" } = {}) {
    if (coupon) await this.page.getByTestId("coupon").fill(coupon);
    await this.page.getByTestId("card").fill(card);
    await this.page.getByTestId("place-order").click();
  }

  errorText() {
    return this.page.getByTestId("checkout-error");
  }
}
