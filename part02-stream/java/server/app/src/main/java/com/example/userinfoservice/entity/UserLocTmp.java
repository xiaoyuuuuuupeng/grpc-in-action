package com.example.userinfoservice.entity;

public class UserLocTmp {
    private String Loc;
    private String Temperature;

    public UserLocTmp(String loc, String temperature) {
        Loc = loc;
        Temperature = temperature;
    }

    public String getLoc() {
        return Loc;
    }

    public void setLoc(String loc) {
        Loc = loc;
    }

    public String getTemperature() {
        return Temperature;
    }

    public void setTemperature(String temperature) {
        Temperature = temperature;
    }
}
