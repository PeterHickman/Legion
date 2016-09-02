# The firewall code is from another project of mine (https://github.com/PeterHickman/Firewall)
# that was developed because ufw didn't work how I needed
# it to. So in the vein of perpetual wheel reinvention I
# went ahead and wrote my own
#
# Some systems will have ufw as the default firewall
# so we need to remove it first

run scripts/package.sh remove ufw

copy scripts/bootstrap/files/fw /root/fw
copy scripts/bootstrap/files/install_fw.sh /root/install_fw.sh
copy scripts/bootstrap/files/iptablesload /root/iptablesload
copy scripts/bootstrap/files/iptablessave /root/iptablessave

ex chmod u+x /root/install_fw.sh
ex ./install_fw.sh
ex fw start
ex fw list

ex rm /root/fw
ex rm /root/install_fw.sh
ex rm /root/iptablesload
ex rm /root/iptablessave
