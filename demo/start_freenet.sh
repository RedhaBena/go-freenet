#!/bin/bash

# Set APP_PATH to the current directory
APP_PATH=$(pwd)

# Define ports and warehouse files as arrays
PORTS=(43210 43211 43212 43213 43214 43215)
WAREHOUSES=("warehouse_a.yaml" "warehouse_b.yaml" "warehouse_c.yaml" "warehouse_d.yaml" "warehouse_e.yaml" "warehouse_f.yaml")

# Loop through ports and warehouses simultaneously
for i in ${!PORTS[@]}; do
  port=${PORTS[$i]}
  warehouse=${WAREHOUSES[$i]}
  
  # Start a new terminal tab (macOS) and run the command in each tab
  osascript <<END
  tell application "Terminal"
    do script "cd $APP_PATH && go run ../. --port $port --warehouse $warehouse"
  end tell
END

done
