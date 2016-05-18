
[web]
Debug = true
Listen = ":1333"
Pprof = ":1777"
Fasthttp = true

[logs.da.console]
Enable = true
Mode = "console"
Level = "Trace"


[logs.default.console]
Enable = true
Mode = "console"
Level = "Trace"

[logs.default.file]
Enable = true
Mode = "file"
Level = "Info"
#model file
FileName  = "./err.log"
LogRotate  = true
MaxLines   = 1000000
MaxSizeShift = 28
DailyRotate  = true
MaxDays     = 7


[logs.default.smtp]
Enable = false
Mode = "smtp"
Level = "Trace"
# model smtp
User      = ""
Passwd    = ""
Host      = ""
Receivers = []
Subject   = ""