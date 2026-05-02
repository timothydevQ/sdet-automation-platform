import { test, expect } from "@playwright/test";
import { LoginPage } from "../pages/LoginPage";
import { createCustomer } from "../fixtures/users";

test.describe("auth @smoke", () => {
  test("user can sign in", async ({ page, request }) => {
    const u = await createCustomer(request);
    const login = new LoginPage(page);
    await login.goto();
    await login.login(u.email, u.password);
    await expect(page).toHaveURL("/", { timeout: 15000 });
  });

  test("invalid credentials show error", async ({ page }) => {
    const login = new LoginPage(page);
    await login.goto();
    await login.login("nobody@nowhere.local", "wrongpassword");
    await expect(login.errorText()).toBeVisible({ timeout: 10000 });
  });
});
