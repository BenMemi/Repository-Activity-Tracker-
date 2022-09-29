#!/bin/bash
touch terraform/start_script.sh
sed -n '5,8p' .env > terraform/start_script.sh
echo 'git clone https://github.com/BenMemi/Repository-Activity-Tracker-.git' >> terraform/start_script.sh
echo 'cd Repository-Activity-Tracker-' >> terraform/start_script.sh
echo 'go build' >> terraform/start_script.sh
echo 'echo "Running tracker"' >> terraform/start_script.sh
echo './tracker' >> terraform/start_script.sh
chmod +x terraform/start_script.sh 


