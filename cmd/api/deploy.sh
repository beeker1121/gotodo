#!/bin/bash
# Build and deploy the app.
#
# This will run the go get and go build commands to
# build the app binary.  This binary is then copied
# to the $APP_DEPLOY_DIR directory with ownership
# and permissions set.
#
# Configure settings so they match the setupvm file
# (if it was used), and make it executable:
#
# chmod u+x deploy.sh
#
# Then execute with sudo:
#
# sudo ./deploy.sh

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
APP_SOURCE_DIR="$ORIG_HOME/work/src/gotodo"

# Go settings
GO_PATH="$ORIG_HOME/work"

# *****************
# Init check
# *****************

# Check if running as root
if [[ $EUID -ne 0 ]]; then
	echo "Not root, please run script with sudo command, exiting..."
	exit 1
fi

# *****************
# Exports
# *****************

echo "Exporting GOPATH and PATH variables..."
export GOPATH="$GO_PATH"
export PATH="$PATH:/usr/local/go/bin:/usr/local/nodejs/bin"

# *****************
# Build and deploy
# *****************

# cd to source directory
echo "cd'ing to $APP_SOURCE_DIR..."
cd $APP_SOURCE_DIR
echo ""

# Copy configuration files to deploy directory
echo "Copying configuration files to $APP_DEPLOY_DIR..."
cp cmd/$APP_BINARY_NAME/config.json $APP_DEPLOY_DIR
echo "Setting ownership and permissions on $APP_DEPLOY_DIR/config.json..."
chown -R $APP_USER:$APP_USER $APP_DEPLOY_DIR/config.json
chmod -R u=r,go= $APP_DEPLOY_DIR/config.json
echo "Permissions set: -r-------- $APP_USER:$APP_USER $APP_DEPLOY_DIR/config.json"
echo ""

# Build binary
echo "cd'ing to cmd/$APP_BINARY_NAME"
cd cmd/$APP_BINARY_NAME
echo "Building $APP_BINARY_NAME binary..."
go build || exit $?
echo "Binary built"
echo ""

# Stop Supervisor
echo "Stopping Supervisor..."
supervisorctl stop $APP_NAME > /dev/null
echo "Supervisor stopped"
echo ""

# Deploy binary
echo "Moving binary to $APP_DEPLOY_DIR..."
mv $APP_BINARY_NAME $APP_DEPLOY_DIR
echo "Setting ownership and permissions on binary..."
chown $APP_USER:$APP_USER $APP_DEPLOY_DIR/$APP_BINARY_NAME
chmod u=rx $APP_DEPLOY_DIR/$APP_BINARY_NAME
chmod go= $APP_DEPLOY_DIR/$APP_BINARY_NAME
echo "Permissions set: -r-x------ $APP_USER:$APP_USER $APP_DEPLOY_DIR/$APP_BINARY_NAME"
echo "Giving CAP_NET_BIND_SERVICE capability to binary..."
setcap "cap_net_bind_service=+ep" $APP_DEPLOY_DIR/$APP_BINARY_NAME
echo ""

# Reload Supervisor
echo "Reloading Supervisor..."
supervisorctl reload > /dev/null
echo "Supervisor reloaded"
echo ""

# *****************
# Finished
# *****************

echo "Deploy complete! Check Supervisor status"
echo ""
echo "sudo supervisorctl status"
echo ""
echo "It should show api as RUNNING"