#!/bin/sh

# We install ruby because some of the scripts that Legion uses are written in ruby
wget https://cache.ruby-lang.org/pub/ruby/2.3/ruby-2.3.1.tar.gz
tar zxf ruby-2.3.1.tar.gz
cd ruby-2.3.1
./configure
make
make install
cd ..
# rm -rf ruby-2.3.1*
