#!/bin/bash
# Install all necessary software on a new VM.
#
# Upload this script to the user home directory,
# configure settings, and make it executable:
#
# chmod u+x setupvm.sh
#
# Then execute with sudo:
#
# sudo ./setupvm.sh
#
# Setup for Debian 8.3 (jessie)
#
# Software Installed:
#
# VIM
# git
# Go
# Node.js
# NPM
# MySQL
# Supervisor
# libcap2-bin (setcap to enable port 80)

# *****************
# Config
# *****************

# User settings
ORIG_USER=$(who -m | awk '{print $1}')
ORIG_HOME=$(getent passwd $ORIG_USER | cut -d: -f6)

# App settings
APP_NAME="gotodoapi"
APP_BINARY_NAME="api"
APP_USER="gotodoapi"
APP_DEPLOY_DIR="/deploy"
APP_LOG_DIR="/var/log/gotodoapi"
APP_LOG_FILE="log.log"

# Download Links
GOLANG_DL="https://dl.google.com/go/go1.10.1.linux-amd64.tar.gz"
GOLANG_DL_NAME_EXT="go1.10.1.linux-amd64.tar.gz"
NODE_DL="https://nodejs.org/dist/v8.10.0/node-v8.10.0-linux-x64.tar.gz"
NODE_DL_NAME="node-v8.10.0-linux-x64"
NODE_DL_NAME_EXT="node-v8.10.0-linux-x64.tar.gz"

# *****************
# Init check
# *****************

# Check if running as root
if [[ $EUID -ne 0 ]]; then
	echo "Not root, please run script with sudo command, exiting..."
	exit 1
fi

# Check if user already exists
if id -u $APP_USER > /dev/null 2>&1; then
	echo "$APP_USER user already exists, exiting..."
	exit 1
fi

# Confirm to proceed
echo "This script will install all necessary software needed and configure"
echo "this VM to run the Go Todo API server.  It should only be used with"
echo "a new VM instance."
echo ""
echo "Checks passed OK:"
echo "  - Running as root"
echo "  - User $APP_USER does not exist"
echo ""
echo "Settings:"
echo "ORIG_USER=$ORIG_USER"
echo "ORIG_HOME=$ORIG_HOME"
echo ""
echo "APP_NAME=$APP_NAME"
echo "APP_BINARY_NAME=$APP_BINARY_NAME"
echo "APP_USER=$APP_USER"
echo "APP_DEPLOY_DIR=$APP_DEPLOY_DIR"
echo "APP_LOG_DIR=$APP_LOG_DIR"
echo "APP_LOG_FILE=$APP_LOG_FILE"
echo ""
echo "GOLANG_DL=$GOLANG_DL"
echo "GOLANG_DL_NAME_EXT=$GOLANG_DL_NAME_EXT"
echo "NODE_DL=$NODE_DL"
echo "NODE_DL_NAME=$NODE_DL_NAME"
echo "NODE_DL_NAME_EXT=$NODE_DL_NAME_EXT"
echo ""

while true; do
	read -p "Do you wish to proceed? Y or n: " prompt
	case $prompt in
		Y)
			break
			;;
		n)
			echo "Exiting..."
			exit 1
			;;
		*)
			echo "Please enter Y or n"
			;;
	esac
done

# *****************
# Create app user
# *****************

echo "Creating new user $APP_USER..."
useradd -M -s /sbin/nologin $APP_USER
echo "$APP_USER user created, no home directory set yet, shell /sbin/nologin (no shell)"


# *****************
# Add backports
# *****************

# Add backports to /etc/apt/sources.list
echo "Adding backports to /etc/apt/sources.list..."
echo "deb http://ftp.debian.org/debian jessie-backports main" >> /etc/apt/sources.list

# *****************
# Update packages
# *****************

# Update packages
echo "Updating packages..."
apt-get -qq update > /dev/null
echo "Packages updated"
echo ""

# *****************
# VIM
# *****************

# Install VIM
echo "Installing VIM..."
apt-get -qq -y install vim > /dev/null
echo "VIM installed"

# Configure VIM
echo "Configuring VIM..."
sed -i 's/"syntax on/syntax on/g' /etc/vim/vimrc
sed -i 's/"set background=dark/set background=dark/g' /etc/vim/vimrc
echo "VIM configured"
echo ""

# *****************
# Git
# *****************

# Install Git
echo "Installing Git..."
apt-get -qq -y install git > /dev/null
echo "Git installed"
echo ""

# *****************
# Go
# *****************

# Install Go
echo "Installing Go binary..."
curl -O $GOLANG_DL
tar -xzf $GOLANG_DL_NAME_EXT -C /usr/local
rm $GOLANG_DL_NAME_EXT
echo "Go installed at /usr/local/go"
mkdir $ORIG_HOME/work
echo "Created workspace for GOPATH at $ORIG_HOME/work"
mkdir $ORIG_HOME/work/bin
mkdir $ORIG_HOME/work/pkg
mkdir $ORIG_HOME/work/src
echo "Created sub directories bin, pkg, src in $ORIG_HOME/work"
chown -R $ORIG_USER:$ORIG_USER $ORIG_HOME/work
echo "Set $ORIG_USER as owner and group of $ORIG_HOME/work recursively"
echo ""

# *****************
# Node.js and npm
# *****************

# Install Node.js and npm
echo "Installing Node.js and npm..."
curl -O $NODE_DL
tar -xzf $NODE_DL_NAME_EXT -C /usr/local
rm $NODE_DL_NAME_EXT
mv /usr/local/$NODE_DL_NAME /usr/local/nodejs
echo "Node.js and npm installed at /usr/local/nodejs"
echo ""

# *****************
# MySQL
# *****************

# Install MySQL
echo "Installing MYSQL..."
apt-get install mysql-server
echo "Running mysql_secure_installation..."
mysql_secure_installation
echo "MySQL installed"
echo ""

# *****************
# App Directory
# *****************

# Set app user home directory
echo "Creating home directory for $APP_USER on disk at $APP_DEPLOY_DIR..."
mkdir -p $APP_DEPLOY_DIR > /dev/null 2>&1
echo "Setting $APP_USER home directory to $APP_DEPLOY_DIR..."
usermod -d $APP_DEPLOY_DIR $APP_USER
echo "Setting ownership and permissions..."
chown -R $APP_USER:$APP_USER $APP_DEPLOY_DIR
chmod -R u=rwX,go= $APP_DEPLOY_DIR
echo "Permissions set: -r-X------ $APP_USER:$APP_USER $APP_DEPLOY_DIR recursively"
echo ""

# *****************
# App Log File
# *****************

# Create app log directory
echo "Creating app log directory at $APP_LOG_DIR..."
mkdir -p $APP_LOG_DIR > /dev/null 2>&1
echo "Creating app log file $APP_LOG_FILE..."
touch $APP_LOG_DIR/$APP_LOG_FILE
echo "Setting ownership and permissions..."
chown -R $APP_USER:$APP_USER $APP_LOG_DIR
chmod -R u=rwX,go= $APP_LOG_DIR
echo "Permissions set: -rw------- $APP_USER:$APP_USER $APP_LOG_FILE"
echo ""

# *****************
# Supervisor
# *****************

# Install Supervisor
echo "Installing Supervisor..."
apt-get -qq -y install supervisor > /dev/null
echo "Supervisor installed"

# Configure Supervisor
echo "Configuring Supervisor..."
echo "Creating supervisor group..."
groupadd --system supervisor
echo "Adding $APP_USER to supervisor group..."
usermod -a -G supervisor $APP_USER

echo "Creating $APP_NAME Supervisor configuration file.."
cat << EOF > /etc/supervisor/conf.d/$APP_NAME.conf
[program:$APP_NAME]
environment=
command=$APP_DEPLOY_DIR/$APP_BINARY_NAME
autostart=true
autorestart=true
startretries=10
user=$APP_USER
directory=$APP_DEPLOY_DIR
redirect_stderr=true
stdout_logfile=/var/log/supervisor/$APP_NAME.log
stdout_logfile_maxbytes=50MB
stdout_logfile_backups=10
EOF

echo "Reloading Supervisor..."
supervisorctl reload > /dev/null
echo "Supervisor configured"
echo ""

# *****************
# libcap2-bin
# *****************

# Install libcap2-bin
echo "Installing libcap2-bin for setcap capabilities..."
apt-get -qq -y install libcap2-bin > /dev/null
echo "libcap2-bin installed"
echo ""

# *********************************************
# /etc/profile
#
# Add new PATH and GOPATH in /etc/profile
# *********************************************

# Add to /etc/profile
echo "Adding new PATH and GOPATH in /etc/profile..."
echo "export PATH=\"\$PATH:/usr/local/go/bin:/usr/local/nodejs/bin\"" >> /etc/profile
echo "export GOPATH=\"$ORIG_HOME/work\"" >> /etc/profile
echo "Added to /etc/profile"
echo ""

# *****************
# Finished
# *****************

echo "Installation complete!"
echo ""
echo "Run source command to export new env variables for Go and Node.js:"
echo ""
echo "source /etc/profile"
echo ""
echo "You can now upload the Go Todo repository to:"
echo "$ORIG_HOME/work/src"
echo ""
echo "Then, modify the /etc/supervisor/conf.d/$APP_NAME.conf file so the correct values"
echo "are set for the needed environment variables, like the database username and password."
echo ""
echo "Make sure to create a new MySQL user and import the cmd/api/schema.sql file to"
echo "your database."
echo ""
echo "Then, cd to the directory, chmod the deploy script, and execute:"
echo ""
echo "cd $ORIG_HOME/work/src/gotodo/cmd/api"
echo "chmod u+x deploy"
echo "sudo ./deploy"
echo ""
echo "The api binary will be built, copied to $APP_DEPLOY_DIR, given permissions"
echo "and setcap capabilities, and then it should be running via Supervisor!"
echo ""
echo "Refer to the README.md file for full installation and deployment instructions."