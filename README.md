# gotodo

Go Todo is an example todo application that implements an HTTP API using Go.

> [!NOTE]
> Please refer to the [DDDStructure](https://github.com/beeker1121/dddstructure) repository for a newer version of this project. In particular, branches in the DDDStructure project show different ways to share top level services.

The goal of this project is to provide a barebones example of how to structure and build a naive, but production ready backend for a web application using Go.

Accompanying this project is an article that goes in-depth on the project structure itself:

**Link coming soon**

Ideally, this project will be a resource for other developers to use and learn from.

It is recommended that you read the above article and study the project in depth before deciding to use it as a base for your own project. Ideally, you would hand-code your own version of this application, using this project and the code herein as a reference.

This readme may seem long-winded and a bit intimidating, but each section is just explaining a very simple action in detail. This is hopefully to provide insight on certain subjects and answer questions that naturally arise.

## Structure Concepts

Concepts involved in using this structure are borrowed from both DDD (Domain Driven Design) and SOA (Service Oriented Architecture).

Mainly, we want to think of our project and the packages contained therein in terms of their responsibilities, or their 'domains'.

For example, the main goal of this project is to provide an HTTP based API. The code for this API will contain our routing and endpoint handlers, middleware, error handling, so on and so forth. Instead of separating this API and its responsibilities into a MVC (Model-View-Controller) architecture or something similar, we have an `api` package with sub-packages handling each of these concerns.

Responding to incoming requests with proper JSON responses is handled by the `api/render` package. Likewise, responding with errors is handled by the `api/errors` package.

We need handlers for each endpoint for version 1 of our API, so those are contained within the `api/v1/handlers` package.

As you can see, the goal is to have each domain manage its own responsibilties.

### Where SOA Gets Involved

We use concepts from SOA to handle the backend actions performed by various API requests.

The interaction with our backend database, MySQL, is handled by the `database` package (or domain). This package has no other concerns other than creating, reading, updating, and deleting entries from the various tables of our database.

The next higher level up is the `services` package. This handles implementing our core application logic using the `database` package to persist and pull what we need.

Finally, our `api` domain uses the `services` package within the various endpoint handlers.

The key is separation of concerns and responsibilties. The `database` package has no concern over what the `api` or `services` package does, its responsibility is simply to handle database interaction. The `api` package has no concern over the database or how it stores data, nor does it care about our application logic; its goal is simply to handle incoming API requests.

## Article Examples

This project will most likely be kept as-is with very little compelxity added.

Within the article, bonuses will be included that will explain how to implement other features using this structure.

These features will hopefully include:

- Member roles.
- Using interfaces to interchange implementations of the `database` and `services` packages.
- TLS using autocert.
- Graceful shutdown.
- Transactions.
- Running on a development vs production server.

## Installation

Download this project from GitHub.

### setupvm.sh and deploy.sh scripts

Within the main application folder, `cmd/api`, there are two files name `setupvm.sh` and `deploy.sh`.

`setupvm.sh` is a bash script that is meant to set up all of the software required for this API server on a new Linux based VM.

`deploy.sh` is a bash script that handles deploying the application.

The `setupvm.sh` script should only be used if you're wanting to install this application on a new VM.

Since this project is meant to show one possible way of running an application in production, we will be using the `deploy.sh` script for deployment.

### Using setupvm.sh

**This script is only meant for brand new VMs. Do not use on an existing VM or computer.**

If you are setting this up on a new VM, you can run the `setupvm.sh` script to install Go, MySQL, and other software.

Upload the `setupvm.sh` script in the `cmd/api` directory to the home directory of your user on your VM, make it executable, and run:

```sh
chmod +x setupvm.sh
sudo ./setupvm.sh
```

### Manual Setup

If you are not using the `setupvm.sh` script, then you will need to manually install some software for our application to run.

1. Install MySQL.

```sh
sudo apt-get install mysql-server
sudo mysql_secure_installation
```

2. Install Supervisor.

```sh
sudo apt-get install supervisor
```

Further configuration of MySQL and Supervisor is explained below.

### Create MySQL Database and Schema

This project requires that you have MySQL installed and running.

A new database in MySQL will be needed to store all of the member and todo data.

It is recommended for security that you create a new user with privileges granted only on this new database.

Create a new database and MySQL user for it, both named `gotodoapi`, replacing `[user_password]` with your own unique password:

```sh
sudo mysql -u root -e "CREATE DATABASE gotodoapi DEFAULT CHARACTER SET utf8mb4 DEFAULT COLLATE utf8mb4_unicode_ci;"
sudo mysql -u root -e "CREATE USER 'gotodoapi'@'localhost' IDENTIFIED BY '[user_password]';"
sudo mysql -u root -e "GRANT ALL ON gotodoapi.* TO 'gotodoapi'@'localhost';"
```

Next, we will need to set up the schema for our database.

Within the `cmd/api` directory, you will see an SQL file named `schema.sql`. This file contains the table creation statements needed to make the `members` and `todos` tables used by our application.

Make sure you have the project uploaded to your Go `src` directory, `cd` to it, then import the schema to our database:

```sh
sudo mysql -u root gotodoapi < cmd/api/schema.sql
```

Awesome! Now our MySQL database is set up and ready for our application.

## Deployment

In testing, it's simple enough to just compile our application to binary and run it manually from the command line.

In production, we will want something that will basically manage running our application for us.

For instance, if we were running our application manually from the command line and our connection closed, our application would close as well.

While there are a number of tools that can handle this for us, one of the easiest and most well known is [Supervisor](http://supervisord.org/).

If you used the `setupvm.sh` script (which should only be used on new servers) to set up your VM, Supervisor is already installed with a configuration file created for the application.

If you did not use the `setupvm.sh` script, you should install Supervisor, and then create the following configuration file for our application at `/etc/supervisor/conf.d/gotodoapi.conf`:

```sh
[program:gotodoapi]
environment=
command=/deploy/api
autostart=true
autorestart=true
startretries=10
user=gotodoapi
directory=/deploy
redirect_stderr=true
stdout_logfile=/var/log/supervisor/gotodoapi.log
stdout_logfile_maxbytes=50MB
stdout_logfile_backups=10
```

You will also notice that for the `user` setting in the configuration file, we have set it to a user named `gotodoapi`.

When running our application in production, we will want to do so as a separate, limited-privileged user. The reason for this is security. If our application becomes compromised, and for example is exploited to allow remote code execution, then that code will be executed as the user running the application. If we were running this application as the root user, it would mean that any code executed by our application would done so as the root user. With root privileges, our entire server would be compromised and under control of the attacker.

On the other hand, if our application is running as a separate user, then code executed by our application will only have the privileges of that user.

If you did not run the `setupvm.sh` script, then you can manually create a new user:

```sh
useradd -M -s /sbin/nologin gotodoapi
```

When running our `deploy.sh` script, the proper permissions will be set allowing this user to run the application.

### Deployment Directory

For deployment, we will want to copy our API binary over to its own directory, separate from the directory containing our source code.

If you used the `setupvm.sh` script, a `deploy` directory has already been created for you at the root level of your drive.

If you did not use the `setupvm.sh` script, create a new root level directory named `deploy`:

```sh
mkdir /deploy
```

### Environment Variables

The first step for deployment is to set the environment variables used to configure our application.

These settings are defined in the `cmd/api/config.json` file, and you can see how they're used in the application via the `cmd/api/main.go` file.

We set sensitive information, such as the MySQL database username and password, using environment variables since we do not want to store this information in the configuration file and hence the project itself. It's possible our repository is shared publicly in one form or another, and we don't want to give out the usernames and passwords to our server and applications.

Since we are using Supervisor to manage running our application, we want to use the Supervisor configuration file for our app to set the environment variables.

Modify the `/etc/supervisor/conf.d/gotodoapi.conf`, adding our environment variables to the `environment` setting:

```sh
[program:gotodoapi]
environment=DB_HOST="localhost",DB_PORT="3306",DB_NAME="gotodoapi",DB_USER="gotodoapi",DB_PASS="",API_HOST="yoururl.com",API_PORT="80",JWT_SECRET="ASecret"
...
```

The `DB_*` fields represent your MySQL database information.

`API_HOST` should be set to the host, or domain, you are running the API on, e.g. `API_HOST="gotodo.io"`

`API_PORT` is the port you wish to run the server on. If you choose port 80, it is recommended you use `setcap` to allow a non-root user to run the application. This is enabled automatically if you are using the `deploy.sh` script for deployment.

`JWT_SECRET` is the string we use to salt our signing key hash for JSON web tokens. This should just be a random string you come up with.

### Running the Deploy Script

So, we now have the following steps completed:

1. MySQL database and schema created with a user named `gotodoapi`.
2. Supervisor installed and configured.
3. A `gotodoapi` OS user created to run our application.
3. A `/deploy` root level directory created to deploy to.
4. Environment variables set to configure our application.

The final step to actually deploy our application is to run the `deploy.sh` script in the `cmd/api` directory.

Upon running this script, it will copy over the `config.json` file to our `/deploy` directory, compile our binary, and copy that over to `/deploy` as well. Permissions are then set on those files so our `gotodapi` user can read and execute them. Supervisor is then reloaded, and our application should be running.

Once you have completed all the above steps, you can run the deploy script from the `cmd/api` directory:

```sh
chmod +x deploy.sh
sudo ./deploy.sh
```

If everything was set up correctly, you should be see the Supervisor status saying `RUNNING` for our application:

```sh
# sudo supervisorctl status
gotodoapi                        RUNNING   pid 333, uptime 12:41:56
```