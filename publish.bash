#!/bin/bash
git commit -am snapshot
git push origin master
git checkout gh-pages
git pull origin master
git commit -am merge
git push origin gh-pages
git checkout master
