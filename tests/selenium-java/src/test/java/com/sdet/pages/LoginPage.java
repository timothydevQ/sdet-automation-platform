package com.sdet.pages;

import org.openqa.selenium.By;
import org.openqa.selenium.WebDriver;

public class LoginPage {
    private final WebDriver driver;

    public LoginPage(WebDriver d) { this.driver = d; }

    public void open(String base) { driver.get(base + "/login"); }

    public void login(String email, String password) {
        driver.findElement(By.cssSelector("[data-testid='email']")).sendKeys(email);
        driver.findElement(By.cssSelector("[data-testid='password']")).sendKeys(password);
        driver.findElement(By.cssSelector("[data-testid='submit']")).click();
    }

    public boolean hasError() {
        return !driver.findElements(By.cssSelector("[data-testid='error']")).isEmpty();
    }
}
