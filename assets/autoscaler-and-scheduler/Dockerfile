# autoscalerdocker/workstation

# docker run -it -v $HOME/workspace:/root/workspace -v $HOME/go:/root/go autoscalerdocker/workstation

FROM fedora:32

MAINTAINER TAS Autoscaler and Scheduler <tas-autoscaler-and-scheduler@groups.vmware.com>

# need ruby to run dns-check.rb & bind-utils for dig & nslookup
RUN dnf update -y; \
  dnf groupinstall -y "Development Tools"; \
  dnf install -y golang; \
  dnf install -y ruby rubygems ruby-devel redhat-rpm-config; \
  dnf install -y mysql-devel openssl-devel zlib-devel; \
  dnf install -y which zsh; \
  dnf install -y npm neovim python3-neovim; \
  gem install yajl-ruby

RUN curl -L https://github.com/cloudfoundry/bosh-cli/releases/download/v6.3.1/bosh-cli-6.3.1-linux-amd64 -o /usr/local/bin/bosh; \
  chmod +x /usr/local/bin/bosh

RUN sh -c "$(curl -fsSL https://raw.github.com/ohmyzsh/ohmyzsh/master/tools/install.sh)"

RUN sh -c "curl vimfiles.luan.sh/install | bash"

# remember to set up MySQL auto/auto_secret and scheduler/password

RUN sed -i 's/robbyrussell/agnoster/' /root/.zshrc

RUN dnf install -y hostname; \
  dnf install -y cmake

CMD [ "/usr/bin/zsh" ]
