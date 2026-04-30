using System.Net.Http.Headers;
using System.Net.Http.Json;

namespace SdetAutomation.Tests.Clients;

public class OrderClient
{
    private readonly HttpClient _http;

    public OrderClient(string baseUrl, string token)
    {
        _http = new HttpClient { BaseAddress = new Uri(baseUrl) };
        _http.DefaultRequestHeaders.Authorization = new AuthenticationHeaderValue("Bearer", token);
    }

    public Task<HttpResponseMessage> AddToCartAsync(string sku, int qty) =>
        _http.PostAsJsonAsync("/cart/items", new { sku, qty });

    public Task<HttpResponseMessage> CheckoutAsync(string coupon = "", string card = "tok_test_visa", string? idemKey = null)
    {
        var req = new HttpRequestMessage(HttpMethod.Post, "/checkout")
        {
            Content = JsonContent.Create(new { coupon, card_token = card }),
        };
        if (idemKey is not null) req.Headers.Add("Idempotency-Key", idemKey);
        return _http.SendAsync(req);
    }

    public Task<HttpResponseMessage> ListAdminOrdersAsync() =>
        _http.GetAsync("/admin/orders");
}
