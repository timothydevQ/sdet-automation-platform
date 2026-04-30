package com.sdet.utils;

import io.github.bonigarcia.wdm.WebDriverManager;
import org.openqa.selenium.WebDriver;
import org.openqa.selenium.chrome.ChromeDriver;
import org.openqa.selenium.chrome.ChromeOptions;

import java.time.Duration;

public class Driver {
    public static WebDriver create() {
        WebDriverManager.chromedriver().setup();
        ChromeOptions opts = new ChromeOptions();
        if (System.getenv("HEADLESS") == null || System.getenv("HEADLESS").equals("true")) {
            opts.addArguments("--headless=new");
        }
        opts.addArguments("--no-sandbox", "--disable-dev-shm-usage", "--window-size=1280,800");
        WebDriver d = new ChromeDriver(opts);
        d.manage().timeouts().implicitlyWait(Duration.ofSeconds(5));
        return d;
    }

    public static String baseUrl() {
        String b = System.getenv("WEB_BASE_URL");
        return b == null ? "http://localhost:3000" : b;
    }

    public static String apiBase() {
        String b = System.getenv("API_BASE");
        return b == null ? "http://localhost:8080" : b;
    }
}
