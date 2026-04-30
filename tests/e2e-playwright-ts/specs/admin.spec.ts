import { test, expect } from "@playwright/test";
import { LoginPage } from "../pages/LoginPage";
import { AdminOrdersPage } from "../pages/AdminOrdersPage";
import { createAdmin, createCustomer } from "../fixtures/users";

test.describe("admin @regression", () => {
  test("admin can view orders", async ({ page, request }) => {
    const a = await createAdmin(request);
    await new LoginPage(page).goto();
    await new LoginPage(page).login(a.email, a.password);
    const orders = new AdminOrdersPage(page);
    await orders.goto();
  });

  test("customer cannot reach admin orders", async ({ page, request }) => {
    const u = await createCustomer(request);
    await new LoginPage(page).goto();
    await new LoginPage(page).login(u.email, u.password);
    await page.goto("/admin/orders");
    await expect(page.getByTestId("admin-error")).toBeVisible();
  });
});
