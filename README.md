#AUTOSTART.SH (ASH)
######The name is misleading, there is very little shell scripting in what is going on here. Sorry about that. The name is just too good to pass up for this.

This is meant as a bootstrapper of client environments. 
####The main goals are ease-of-use for:
- initial configuration 
- installation and usage 
- maintenance

####Specific non-goals:
- Speed (Not like it's really slow either.)
- Efficiency (no state tracking, that's the package manager's job)

For you docker nerds:
Think of this as a Dockerfile for... you. It builds your environment for you.

For everyone else, it's a combination of a dotfile manager and a very basic universal package manager. Since it can potentially support nearly any *nix style environment (and maybe even Windows), that means it is the lowest common denominator of all of these environments. Not to speak ill of my own sofware, but that means it's very limited and not very smart. However, it makes creating the same set of personal environments to be applied in many different contexts/operating systems very easy and fast. Not to mention it's just awesome.

This is heavily inspired by other dotfile and home directory managers. Specifically, homely and homemaker.

## Installation
TODO: add installation line. The intent is a single curl call that both downloads the binary for the appropriate architecture/environment, but also run a personal configuration to start the ASH process
```bash
curl <someUrl> - runs script - dls and installs binary - optionally also runs passed in configuration
```

### As a library
This package was designed with the intent of being consumed as a library. You can use it in go by:
```bash
go get github.com/morganhein/autostart.sh
```

## Overview of how this works
The intent is for an individual to write and maintain a single file that declares all the steps required to set up and bootstrap an environment in a specific way. If the user desires, they can also use this to manage their configuration and misc dotfiles.

The originally supported environments are: 
- arch linux, using pacman and yay
- fedora and friends, using dnf
- alpine linux, using apk
- ubuntu and friends, using apt
- mac using brew
- special `sh` and `bash` targets for downloading and running scripts not available via package managers

This adds another attack vector for you, the user, because now there is automation in your setup pipeline. You must audit both
this software, and any scripts/commands you want run during setup, to ensure they meet your security requirements. This software, by itself, 
does nothing. However, given a malicious configuration file, it is very possible for you to install something with a malicious intent.

More information on the security posture of ASH can be found below, under security. (TODO: LINK MAYBE?)

Once a user writes a configuration file, they can apply that config to their current environment. It will attempt to install the software and link any configuration/dotfiles that the user requested.

Once a system is up and running, there are facilities for maintaining and updating the dotfiles and ASH configuration.

### Usage:

need functions for:
- running config(s)
    - interactive mode
    - using a prioritized list of installer targets
    - run without applying links
    - run without installing packages
- updating links and dotfiles

`ash run [configuration file] --installers=apk,sh`

#### Available environment variables available in cmd lines
- original_task - the name of the original requested task
- current_task - the name of the current task being run
- url_file - the filepath of the downloaded file, as specified in the installerOpts.url, or the result of the @download macro
- dl_file - the result of the @download <path> macro
- pkg - pkg name
- installer - the name of the installer being used
- sudo - inserts sudo if enabled
- last_cmd - the full text of the last cmd run
- next_cmd - the full text of the next cmd to run (after parsing all values)
- link_dest - the link destination for link creation
- link_src - the source directory containing original files to link to
- config_path - the path to the config file
