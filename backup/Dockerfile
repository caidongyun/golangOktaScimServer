FROM ubuntu

#Add Okta2Anything
ADD https://dev-337217-admin.okta.com/static/ldap-agent/OktaLDAPAgent-05.05.07_amd64.deb /OktaLDAPAgentProd.deb
ADD okta2anything.linux /
ADD promiscuous.js /

RUN apt-get update -y
#RUN apt-get install -y curl npm wget vim python-setuptools python-dev python-pip build-essential git make ssh zip
RUN apt-get install -y curl npm wget vim build-essential git make ssh zip
RUN rm -rf /var/lib/apt/lists/*

#install pip
#RUN easy_install pip

#RUN apt-get install python-setuptools python-dev curl vim npm git make ssh zip -y
RUN npm config set registry http://registry.npmjs.org/

# Install node.js
RUN curl --silent --location https://deb.nodesource.com/setup_10.x | bash -
RUN apt-get install --yes nodejs


# Bundle app source
# Trouble with COPY http://stackoverflow.com/a/30405787/2926832
COPY . /src

#add Okta LDAP Agent

# Install app dependencies
RUN cd /src; npm install



# Adding Dockertest components







