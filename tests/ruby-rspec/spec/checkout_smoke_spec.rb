require_relative "support/api"

RSpec.describe "Checkout smoke", :smoke do
  it "places an order from cart" do
    email, pw = Api.register_customer
    visit "/login"
    fill_in "Email", with: email
    fill_in "Password", with: pw
    click_button "Login"

    expect(page).to have_css("[data-testid^='product-']", wait: 10)
    first("[data-testid^='add-']").click

    visit "/cart"
    click_button "Checkout"

    fill_in "Card token", with: "tok_test_visa"
    click_button "Place order"

    expect(page).to have_current_path(%r{/orders/}, wait: 10)
  end
end
