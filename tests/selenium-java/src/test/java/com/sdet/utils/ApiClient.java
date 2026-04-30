package com.sdet.utils;

import java.net.URI;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;
import java.util.UUID;

public class ApiClient {
    private static final HttpClient client = HttpClient.newHttpClient();

    public static String registerCustomer() throws Exception {
        String email = "u" + UUID.randomUUID().toString().replace("-", "").substring(0, 10) + "@test.local";
        return register(email, "Hunter22!");
    }

    public static String registerAdmin() throws Exception {
        String email = "a" + UUID.randomUUID().toString().replace("-", "").substring(0, 10) + "@admin.local";
        return register(email, "Hunter22!");
    }

    private static String register(String email, String pw) throws Exception {
        String body = String.format("{\"email\":\"%s\",\"password\":\"%s\"}", email, pw);
        HttpRequest req = HttpRequest.newBuilder()
                .uri(URI.create(Driver.apiBase() + "/auth/register"))
                .header("Content-Type", "application/json")
                .POST(HttpRequest.BodyPublishers.ofString(body))
                .build();
        HttpResponse<String> r = client.send(req, HttpResponse.BodyHandlers.ofString());
        if (r.statusCode() / 100 != 2) throw new RuntimeException("register: " + r.body());
        return email + "|Hunter22!";
    }
}
