using System.Net;
using NUnit.Framework;
using SdetAutomation.Tests.Clients;

namespace SdetAutomation.Tests.Tests;

[TestFixture]
[Category("smoke")]
public class PaymentValidationTests
{
    private static readonly string ApiBase =
        Environment.GetEnvironmentVariable("API_BASE") ?? "http://localhost:8080";

    [Test]
    public async Task DeclinedCardReturns402()
    {
        var auth = new AuthClient(ApiBase);
        var token = await auth.RegisterCustomerAsync();
        var client = new OrderClient(ApiBase, token.Token);

        await client.AddToCartAsync("SKU-001", 1);
        var resp = await client.CheckoutAsync(card: "tok_decline_card");

        Assert.That((int)resp.StatusCode, Is.EqualTo(402));
    }

    [Test]
    public async Task SuccessfulCardReturns201()
    {
        var auth = new AuthClient(ApiBase);
        var token = await auth.RegisterCustomerAsync();
        var client = new OrderClient(ApiBase, token.Token);

        await client.AddToCartAsync("SKU-001", 1);
        var resp = await client.CheckoutAsync();

        Assert.That(resp.StatusCode, Is.EqualTo(HttpStatusCode.Created));
    }
}
