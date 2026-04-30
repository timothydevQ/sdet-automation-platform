module Api
  def self.register_customer
    email = "u#{SecureRandom.hex(5)}@test.local"
    password = "Hunter22!"
    HTTParty.post("#{API_BASE}/auth/register",
      body: { email: email, password: password }.to_json,
      headers: { "Content-Type" => "application/json" })
    [email, password]
  end

  def self.register_admin
    email = "a#{SecureRandom.hex(5)}@admin.local"
    password = "Hunter22!"
    HTTParty.post("#{API_BASE}/auth/register",
      body: { email: email, password: password }.to_json,
      headers: { "Content-Type" => "application/json" })
    [email, password]
  end
end
