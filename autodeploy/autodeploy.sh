#!/bin/bash
cd /home/ayang/site/Ayang-Blog
git pull origin main
git submodule update --remote themes/FixIt # 更新子项目
/home/ayang/go/bin/hugo
echo $(date "+%Y-%m-%d %H:%M:%S") >> autodeploy/autodeploy.log
