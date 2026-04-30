require_relative "support/api"

RSpec.describe "Admin smoke", :regression do
  it "blocks customers from admin orders" do
    email, pw = Api.register_customer
    visit "/login"
    fill_in "Email", with: email
    fill_in "Password", with: pw
    click_button "Login"

    visit "/admin/orders"
    expect(page).to have_css("[data-testid='admin-error']")
  end
end
