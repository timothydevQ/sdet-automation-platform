using System;
using System.Net.Http;
using System.Net.Http.Json;
using System.Threading.Tasks;

namespace SdetAutomation.Tests.Clients;

public record TokenResponse(string Token, string Role);

public class AuthClient
{
    private readonly HttpClient _http;

    public AuthClient(string baseUrl)
    {
        _http = new HttpClient { BaseAddress = new Uri(baseUrl) };
    }

    public async Task<TokenResponse> RegisterCustomerAsync()
    {
        var email = "u" + Guid.NewGuid().ToString("N").Substring(0, 10) + "@test.local";
        return await RegisterAsync(email, "Hunter22!");
    }

    public async Task<TokenResponse> RegisterAdminAsync()
    {
        var email = "a" + Guid.NewGuid().ToString("N").Substring(0, 10) + "@admin.local";
        return await RegisterAsync(email, "Hunter22!");
    }

    private async Task<TokenResponse> RegisterAsync(string email, string password)
    {
        var resp = await _http.PostAsJsonAsync("/auth/register", new { email, password });
        resp.EnsureSuccessStatusCode();
        return (await resp.Content.ReadFromJsonAsync<TokenResponse>())!;
    }
}
