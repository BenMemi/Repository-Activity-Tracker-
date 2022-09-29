#!/bin/bash
touch terraform/start_script.sh
sed -n '5,8p' .env > terraform/start_script.sh
echo 'sudo wget -c https://go.dev/dl/go1.19.1.linux-amd64.tar.gz' >> terraform/start_script.sh
echo 'sudo tar -C /usr/local -xvzf go1.19.1.linux-amd64.tar.gz' >> terraform/start_script.sh
echo 'export  PATH=$PATH:/usr/local/go/bin' >> terraform/start_script.sh
echo 'git clone https://github.com/BenMemi/Repository-Activity-Tracker-.git' >> terraform/start_script.sh
echo 'cd Repository-Activity-Tracker-' >> terraform/start_script.sh
echo 'go build' >> terraform/start_script.sh
echo 'echo "Running tracker"' >> terraform/start_script.sh
echo 'chmod 777 terraform/start_script.sh' >> terraform/start_script.sh
echo './tracker' >> terraform/start_script.sh




