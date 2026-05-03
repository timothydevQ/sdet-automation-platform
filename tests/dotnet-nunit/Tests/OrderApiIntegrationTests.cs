using System;
using System.Threading.Tasks;
using NUnit.Framework;
using SdetAutomation.Tests.Clients;

namespace SdetAutomation.Tests.Tests;

[TestFixture]
[Category("regression")]
public class OrderApiIntegrationTests
{
    private static readonly string ApiBase =
        Environment.GetEnvironmentVariable("API_BASE") ?? "http://localhost:8080";

    [Test]
    public async Task IdempotencyKeyReturnsSameOrder()
    {
        var auth = new AuthClient(ApiBase);
        var token = await auth.RegisterCustomerAsync();
        var client = new OrderClient(ApiBase, token.Token);

        await client.AddToCartAsync("SKU-001", 1);
        var key = Guid.NewGuid().ToString("N");

        var first = await client.CheckoutAsync(idemKey: key);
        Assert.That((int)first.StatusCode, Is.EqualTo(201));

        await client.AddToCartAsync("SKU-001", 1);
        var second = await client.CheckoutAsync(idemKey: key);
        Assert.That((int)second.StatusCode, Is.EqualTo(200));
    }

    [Test]
    public async Task CustomerCannotListAdminOrders()
    {
        var auth = new AuthClient(ApiBase);
        var token = await auth.RegisterCustomerAsync();
        var client = new OrderClient(ApiBase, token.Token);

        var resp = await client.ListAdminOrdersAsync();
        Assert.That((int)resp.StatusCode, Is.EqualTo(403));
    }
}
