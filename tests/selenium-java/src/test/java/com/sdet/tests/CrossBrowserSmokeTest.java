package com.sdet.tests;

import com.sdet.pages.CatalogPage;
import com.sdet.utils.Driver;
import org.junit.jupiter.api.AfterEach;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Tag;
import org.junit.jupiter.api.Test;
import org.openqa.selenium.WebDriver;

import static org.junit.jupiter.api.Assertions.assertTrue;

@Tag("regression")
public class CrossBrowserSmokeTest {
    private WebDriver driver;

    @BeforeEach
    void setup() { driver = Driver.create(); }

    @AfterEach
    void tear() { if (driver != null) driver.quit(); }

    @Test
    void catalogLoads() {
        CatalogPage cat = new CatalogPage(driver);
        cat.open(Driver.baseUrl());
        cat.waitForProducts();
        assertTrue(cat.isLoaded());
    }
}
