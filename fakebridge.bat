@echo off

:: Set environment variable
set MNEMONIC=all all all all all all all all all all all all
set PASSWORD=

:: Stop the "trezord" process if it is running
taskkill /F /IM trezord.exe 2>NUL

:: Run the fakebridge command with index 0
fakebridge.exe -index 0

:: Keep the window open to display any error messages
pause