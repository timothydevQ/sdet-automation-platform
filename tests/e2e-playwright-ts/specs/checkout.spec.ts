import { test, expect } from "@playwright/test";
import { LoginPage } from "../pages/LoginPage";
import { CatalogPage } from "../pages/CatalogPage";
import { CartPage } from "../pages/CartPage";
import { CheckoutPage } from "../pages/CheckoutPage";
import { createCustomer } from "../fixtures/users";

test.describe("checkout @smoke", () => {
  test("end to end happy path", async ({ page, request }) => {
    const u = await createCustomer(request);
    await new LoginPage(page).goto();
    await new LoginPage(page).login(u.email, u.password);

    const catalog = new CatalogPage(page);
    await catalog.goto();
    await catalog.waitForProducts();
    await catalog.addToCart("SKU-001");

    const cart = new CartPage(page);
    await cart.goto();
    await cart.expectItemCount(1);
    await cart.checkout();

    await new CheckoutPage(page).submit({ coupon: "WELCOME10" });
    await expect(page).toHaveURL(/\/orders\//);
  });

  test("declined card surfaces error", async ({ page, request }) => {
    const u = await createCustomer(request);
    await new LoginPage(page).goto();
    await new LoginPage(page).login(u.email, u.password);
    const catalog = new CatalogPage(page);
    await catalog.goto();
    await catalog.waitForProducts();
    await catalog.addToCart("SKU-001");
    await new CartPage(page).goto();
    await new CartPage(page).checkout();
    const checkout = new CheckoutPage(page);
    await checkout.submit({ card: "tok_decline_card" });
    await expect(checkout.errorText()).toBeVisible();
  });
});
