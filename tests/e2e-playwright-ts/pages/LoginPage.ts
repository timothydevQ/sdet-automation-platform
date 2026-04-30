import { Page, expect } from "@playwright/test";

export class LoginPage {
  constructor(private page: Page) {}

  async goto() {
    await this.page.goto("/login");
    await expect(this.page.getByTestId("login-page")).toBeVisible();
  }

  async login(email: string, password: string) {
    await this.page.getByTestId("email").fill(email);
    await this.page.getByTestId("password").fill(password);
    await this.page.getByTestId("submit").click();
  }

  errorText() {
    return this.page.getByTestId("error");
  }
}
