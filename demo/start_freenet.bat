@echo off

:: Set APP_PATH to the current directory
set APP_PATH=%cd%

:: Define ports and warehouse files
set PORTS=43210 43211 43212 43213 43214 43215
set WAREHOUSES=warehouse_a.yaml warehouse_b.yaml warehouse_c.yaml warehouse_d.yaml warehouse_e.yaml warehouse_f.yaml

setlocal enabledelayedexpansion

:: Convert PORTS and WAREHOUSES into indexed variables
set "ports[0]=43210"
set "ports[1]=43211"
set "ports[2]=43212"
set "ports[3]=43213"
set "ports[4]=43214"
set "ports[5]=43215"

set "warehouses[0]=warehouse_a.yaml"
set "warehouses[1]=warehouse_b.yaml"
set "warehouses[2]=warehouse_c.yaml"
set "warehouses[3]=warehouse_d.yaml"
set "warehouses[4]=warehouse_e.yaml"
set "warehouses[5]=warehouse_f.yaml"

:: Loop through ports and warehouses simultaneously
for /L %%i in (0,1,5) do (
    set "port=!ports[%%i]!"
    set "warehouse=!warehouses[%%i]!"
    start cmd /k "cd /d %APP_PATH% && go run ../. --port !port! --warehouse !warehouse!"
)

