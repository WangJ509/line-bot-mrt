# Line Bot MRT

This project is a line bot that helps you to retrive the information from [MOTC Transport API](https://ptx.transportdata.tw/MOTC/?urls.primaryName=%E8%BB%8C%E9%81%93V2#/Metro/MetroApi_Frequency_2100). As it is in the early stage of develepment, this project only support the time table function.

I have deployed this project on Heroku, if you want to have a look, you can join this official account via the following link.

![](https://i.imgur.com/awb3Jfq.png)

## Motivation

Recently, Taipei MRT company kept promote their new app `台北捷運go`. However, I found the app is full ADs and not even available for checking the time table!

As a result, I decided to find a more convenient way to get the time table information. I found [MOTC Transport API](https://ptx.transportdata.tw/MOTC/?urls.primaryName=%E8%BB%8C%E9%81%93V2#/Metro/MetroApi_Frequency_2100) and decided to combine with line bot.

I choose line bot because Line app supports Shortcut App on my iphone, so I can write a routine to automatically check the time table before I go to work.

## Usage
To retrieve the time table of Taipei MRT, type
```
時刻表 出發站 終點站 數量
```
For example, send `時刻表 景美 松山 3` will response with 
```
21:44
21:52
22:01
```

If the message does not start with a valid query command, it will response with `りしれ供さ小`


## How to run the program
First install `go` with version `1.18`
To install the dependency, run 
```
$ go mod tidy
```
Before executing the program, you should set the following three `environment variable`
- `ChannelSecret`: you can get the channel secret from line developer website
- `ChannelAccessToken`: you can also get this token from line developer website
- `PORT`: the port of this API service.

After setting the environment variables, you can execute this program by
```
$ go run main.go
```
To check the program is actually running, you can visit 
```
http://localhost:[your_port]/health-check
```

To test the program, you can run 
```
$ go test -v [package]
```
For example, `go test -v ./mrt` will do the test I wrote in `mrt` package.

## Future Work
- Add support for suggesting station name. Like `Did you mean: xxx` in Google search.
- Save data in DB to reduce the time for querying MOTC API.
- Add support for checking national holidays, as the time table would be different for national holidays.
- Add routing check, so that user can input any destination, and the program should check how to travel to the destination.
- Add other operations, e.g. real-time information of MRT from MOTC API.
