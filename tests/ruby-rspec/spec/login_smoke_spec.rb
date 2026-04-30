require_relative "support/api"

RSpec.describe "Login smoke", :smoke do
  it "logs in a registered user" do
    email, pw = Api.register_customer
    visit "/login"
    fill_in "Email", with: email
    fill_in "Password", with: pw
    click_button "Login"
    expect(page).to have_css("[data-testid='catalog-page']")
  end

  it "shows error on bad credentials" do
    visit "/login"
    fill_in "Email", with: "noone@test.local"
    fill_in "Password", with: "wrong"
    click_button "Login"
    expect(page).to have_css("[data-testid='error']")
  end
end
