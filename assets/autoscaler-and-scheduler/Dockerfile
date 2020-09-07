# autoscalerdocker/workstation

# docker run -it -v $HOME/workspace:/root/workspace -v $HOME/go:/root/go autoscalerdocker/workstation

FROM fedora:32

MAINTAINER TAS Autoscaler and Scheduler <tas-autoscaler-and-scheduler@groups.vmware.com>

RUN dnf update -y; \
  dnf groupinstall -y "Development Tools"; \
  dnf install -y golang; \
  dnf install -y ruby rubygems ruby-devel redhat-rpm-config; \
  dnf install -y mysql-devel openssl-devel zlib-devel; \
  dnf install -y which zsh; \
  dnf install -y curl npm ripgrep fd-find neovim python3-neovim; \
  dnf install -y cmake direnv hostname; \
  dnf install -y procps-ng tcpdump lsof net-tools iproute bind-utils; \
  gem install yajl-ruby

RUN curl -L https://github.com/cloudfoundry/bosh-cli/releases/download/v6.3.1/bosh-cli-6.3.1-linux-amd64 -o /usr/local/bin/bosh; \
  chmod +x /usr/local/bin/bosh

RUN sh -c "$(curl -fsSL https://raw.github.com/ohmyzsh/ohmyzsh/master/tools/install.sh)"

RUN git clone https://github.com/luan/nvim /root/.config/nvim

# remember to set up MySQL auto/auto_secret and scheduler/password

RUN sed -i 's/robbyrussell/agnoster/' /root/.zshrc

RUN curl -L https://github.com/clvv/fasd/tarball/1.0.1 -o - | tar xzvf - ; \
  cd clvv*; \
  make install

CMD [ "/usr/bin/zsh" ]
