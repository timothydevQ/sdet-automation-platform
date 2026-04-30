package com.sdet.pages;

import org.openqa.selenium.By;
import org.openqa.selenium.WebDriver;
import org.openqa.selenium.support.ui.ExpectedConditions;
import org.openqa.selenium.support.ui.WebDriverWait;

import java.time.Duration;

public class CatalogPage {
    private final WebDriver driver;

    public CatalogPage(WebDriver d) { this.driver = d; }

    public void open(String base) { driver.get(base + "/"); }

    public void waitForProducts() {
        new WebDriverWait(driver, Duration.ofSeconds(10))
                .until(ExpectedConditions.presenceOfElementLocated(
                        By.cssSelector("[data-testid^='product-']")));
    }

    public boolean isLoaded() {
        return !driver.findElements(By.cssSelector("[data-testid='catalog-page']")).isEmpty();
    }
}
