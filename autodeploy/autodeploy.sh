#!/bin/bash
cd /home/ayang/site/Ayang-Blog
git pull origin main
git submodule update --remote themes/FixIt # 更新子项目
/home/ayang/go/bin/hugo  # 执行hugo构建静态网站
cp /home/ayang/site/Ayang-Blog/public/posts/index.html /home/ayang/site/Ayang-Blog/public/index.html # 把 posts 作为主页
echo $(date "+%Y-%m-%d %H:%M:%S") >> autodeploy/autodeploy.log
