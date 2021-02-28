# Deployer Agent

This repository contains the source code for the Deployer Agent. A small Go based package that is used to synchronize your account SSH keys. Please note that to use this you'll need to first sign up for an account at https://deployer.codions.com

In most cases you can use the installer provided inside your Deployer account, which obtains the latest stable release of the agent from Github. However if you'd prefer to build the package from source, please see below.

## System Requirements

Whilst the Deployer Agent is a fairly small and minimal package, it does have a few requirements. The first of which is our list of supported operating systems.

Whilst the agent is likely to work fine on other *nix based systems, these are the ones we currently officially support:

- CentOS
    - CentOS 8
    - CentOS 7
    - CentOS 6
- Debian
    - Debian 10 (Buster)
    - Debian 9 (Stretch)
    - Debian 8 (Jessie)
    - Debian 7 (Wheezy)

- Fedora
    - Fedora 32
    - Fedora 31
    - Fedora 30
    - Fedora 29
    - Fedora 28

- Ubuntu
    - Ubuntu 20.04
    - Ubuntu 18.04
    - Ubuntu 16.04

## Permissions

The Deployer Agent must be run as a high-level user, with permission to modify files owned by another user.

It is assumed that the agent will be running as the `root` user, however if you are running as another user and have allocated the correct passwordless sudo permissions then you can modify the system cron job, or can manually trigger the `deployer sync` command.


## Available Commands

The agent includes a number of commands. These include the ability to add a new system account, remove an existing system account, and trigger a manual sync of all accounts.

Details on each command can be returned by running `deployer --help` from command line.

## Changelog

Please see [Releases](../../releases) for more information what has changed recently.

## Security Vulnerabilities

If you discover any security related issues, please email dev@codions.com instead of using the issue tracker.

## License

The Laravel framework is open-sourced software licensed under the [MIT license](https://opensource.org/licenses/MIT).

