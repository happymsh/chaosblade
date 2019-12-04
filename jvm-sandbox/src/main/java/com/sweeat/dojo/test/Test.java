package com.sweeat.dojo.test;
import java.util.concurrent.TimeUnit;
public class Test {

    public String test() {
        return "test...";
    }

    public static void main(String[] args) throws InterruptedException {
        Test testJVM = new Test();

        while (true) {
            try {
                System.out.println(testJVM.test());
            } catch (Exception e) {
                System.out.println(e.getMessage());
            }
            TimeUnit.SECONDS.sleep(1);
        }
    }
}