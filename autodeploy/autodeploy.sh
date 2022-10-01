#!/bin/bash
cd /home/ayang/site/Ayang-Blog
git pull
/home/ayang/go/bin/hugo
echo $(date "+%Y-%m-%d %H:%M:%S") >> autodeploy/autodeploy.log
