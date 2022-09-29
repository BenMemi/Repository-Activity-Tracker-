#!/bin/bash
touch terraform/start_script.sh
sed -n '1,100p' .env > terraform/start_script.sh



