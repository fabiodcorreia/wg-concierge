# -*- mode: ruby -*-
# vi: set ft=ruby :

BOX_NAME="wg-concierge-dev-box"
BOX_CORES="2"
BOX_MEM="1024"

Vagrant.configure("2") do |config|

  config.vm.box = "debian/buster64"
  config.vm.hostname = "#{BOX_NAME}"
  config.vm.synced_folder "./build", "/home/vagrant/wg-concierge", type: "rsync", rsync__verbose: "true"

  config.vm.provider "virtualbox" do |vb|
    vb.name = "#{BOX_NAME}"
    vb.cpus = "#{BOX_CORES}"
    vb.memory = "#{BOX_MEM}"
    vb.gui = false
  end

  config.vm.provision "shell", inline: <<-SHELL
    echo "deb http://deb.debian.org/debian buster-backports main" | tee --append /etc/apt/sources.list
    sh -c 'printf "Package: *\nPin: release a=unstable\nPin-Priority: 90\n" > /etc/apt/preferences.d/limit-unstable'
    apt-get update
    apt-get upgrade -y
    apt-get install wireguard -y

    echo "[Interface]\nAddress = 10.1.1.1/24\nSaveConfig = true\nListenPort = 55555\nPrivateKey = qJvFeHHuffBaPWx4veJGQqXw6j5zdo5cSOaBd1Z0Km4=" | tee /etc/wireguard/wg0.conf

    sysctl -w net.ipv4.ip_forward=1
    systemctl enable wg-quick@wg0

  SHELL
end

