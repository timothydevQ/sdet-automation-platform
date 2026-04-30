package com.sdet.tests;

import com.sdet.pages.CatalogPage;
import com.sdet.pages.LoginPage;
import com.sdet.utils.ApiClient;
import com.sdet.utils.Driver;
import org.junit.jupiter.api.AfterEach;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.Tag;
import org.openqa.selenium.WebDriver;

import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertTrue;

@Tag("smoke")
public class LoginSmokeTest {
    private WebDriver driver;

    @BeforeEach
    void setup() { driver = Driver.create(); }

    @AfterEach
    void tear() { if (driver != null) driver.quit(); }

    @Test
    void userCanLogin() throws Exception {
        String creds = ApiClient.registerCustomer();
        String[] parts = creds.split("\\|");
        LoginPage login = new LoginPage(driver);
        login.open(Driver.baseUrl());
        login.login(parts[0], parts[1]);
        CatalogPage catalog = new CatalogPage(driver);
        catalog.waitForProducts();
        assertTrue(catalog.isLoaded());
    }

    @Test
    void invalidCredentialsShowError() {
        LoginPage login = new LoginPage(driver);
        login.open(Driver.baseUrl());
        login.login("nobody@test.local", "bad");
        assertTrue(login.hasError());
    }

    @Test
    void emptyFormDoesNotSubmit() {
        LoginPage login = new LoginPage(driver);
        login.open(Driver.baseUrl());
        login.login("", "");
        assertFalse(driver.getCurrentUrl().endsWith("/"));
    }
}
