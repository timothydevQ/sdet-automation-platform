package com.sdet.pages;

import org.openqa.selenium.By;
import org.openqa.selenium.WebDriver;
import org.openqa.selenium.WebElement;
import org.openqa.selenium.support.ui.ExpectedConditions;
import org.openqa.selenium.support.ui.WebDriverWait;

import java.time.Duration;

public class LoginPage {
    private final WebDriver driver;
    private final WebDriverWait wait;

    public LoginPage(WebDriver d) {
        this.driver = d;
        this.wait = new WebDriverWait(d, Duration.ofSeconds(15));
    }

    public void open(String base) {
        driver.get(base + "/login");
        wait.until(ExpectedConditions.presenceOfElementLocated(
                By.cssSelector("[data-testid='email']")));
    }

    public void login(String email, String password) {
        driver.findElement(By.cssSelector("[data-testid='email']")).sendKeys(email);
        driver.findElement(By.cssSelector("[data-testid='password']")).sendKeys(password);
        driver.findElement(By.cssSelector("[data-testid='submit']")).click();
    }

    public boolean hasError() {
        try {
            WebElement err = wait.until(ExpectedConditions.visibilityOfElementLocated(
                    By.cssSelector("[data-testid='error']")));
            return err.isDisplayed();
        } catch (Exception e) {
            return false;
        }
    }
}