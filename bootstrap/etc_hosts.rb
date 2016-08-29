#!/usr/bin/env ruby
# encoding: UTF-8

HOST_FILE = '/etc/hosts'

class IPAddress
  def self.to_i(address)
    x = address.split('.').map(&:to_i)
    t = 0
    while x.any?
      y = x.shift
      t = (t * 256) + y
    end
    t
  end

  def self.private?(address)
    x = to_i(address)

    if 167_772_161 <= x && x <= 184_549_374
      # ["10.0.0.0/8", "10.0.0.1", 167772161, "10.255.255.254", 184549374]
      true
    elsif 2_886_729_729 <= x && x <= 2_887_778_302
      # ["172.16.0.0/12", "172.16.0.1", 2886729729, "172.31.255.254", 2887778302]
      true
    elsif 3_232_235_521 <= x && x <= 3_232_301_054
      # ["192.168.0.0/16", "192.168.0.1", 3232235521, "192.168.255.254", 3232301054]
      true
    else
      false
    end
  end

  def self.special?(address)
    if address == '127.0.0.1'
      true
    elsif address == '127.0.1.1'
      true
    elsif address == '255.255.255.255'
      true
    else
      false
    end
  end

  def self.address?(address)
    x = address.split('.').map(&:to_i)
    if x.size == 4
      (0..3).each do |i|
        return false if x[i] < 0 || x[i] > 255
      end
      true
    else
      false
    end
  end
end

def host_ip_address
  t = ['127.0.1.1']

  %x(ifconfig).split("\n").each do |line|
    next unless line.include?('inet addr:')

    addresses = line.split(/\s+/).select { |y| y.index('addr:') }.map { |y| y[5..-1] }

    addresses.each do |address|
      if IPAddress.address?(address)
        t << address unless IPAddress.special?(address)
      end
    end
  end

  t.uniq
end

hostname = %x(hostname).chomp

host_addresses = host_ip_address

t = ''
changes = false

File.open(HOST_FILE, 'r').each do |line|
  x = line.split(/\s+/)
  if x.size > 1
    address = x.shift

    if IPAddress.address?(address)
      if host_addresses.include?(address)
        y = "#{address}\t#{hostname}\n"
        changes = true if y != line
        t << y
      else
        t << line
      end
    else
      t << line
    end
  else
    t << line
  end
end

if changes
  puts "Updated #{HOST_FILE} file:"
  puts
  puts t
  f = File.open(HOST_FILE, 'w')
  f.puts t
  f.close
else
  puts "No changes for the #{HOST_FILE} files"
end
