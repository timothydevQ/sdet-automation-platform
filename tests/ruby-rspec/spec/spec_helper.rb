require "capybara/rspec"
require "selenium-webdriver"
require "httparty"
require "securerandom"

API_BASE = ENV.fetch("API_BASE", "http://localhost:8080")
WEB_BASE = ENV.fetch("WEB_BASE_URL", "http://localhost:3000")

Capybara.app_host = WEB_BASE
Capybara.run_server = false
Capybara.default_max_wait_time = 10

Capybara.register_driver :headless_chrome do |app|
  options = Selenium::WebDriver::Chrome::Options.new
  options.add_argument("--headless=new")
  options.add_argument("--no-sandbox")
  options.add_argument("--disable-dev-shm-usage")
  options.add_argument("--window-size=1280,800")
  Capybara::Selenium::Driver.new(app, browser: :chrome, options: options)
end

Capybara.default_driver = :headless_chrome
Capybara.javascript_driver = :headless_chrome

RSpec.configure do |config|
  config.expect_with :rspec do |c|
    c.syntax = :expect
  end
end
