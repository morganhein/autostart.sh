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
    - run without applying links
    - run without installing packages
- updating links and dotfiles

`ash run [configuration file] --installers=gvm,brew <task>`

#### Available environment variables available in cmd lines
- sudo: if sudo should be enabled for commands
- pkg: pkg name
- installer: the name of the installer being used
- sudo: inserts sudo if enabled
- link_dest: the link destination for link creation
- link_src: the source directory containing original files to link to
- config_path: the path to the config file

## Built-In Macros

The only macros available are 
`@install` and `@download` macros.
The install macro can be overridden for each platform, but the download macro uses a go function, so cannot be overridden.

## Task format

```yaml
[task.norman] # installs the norman keyboard layout
    cmds = ["${sudo} cp /etc/default/keyboard /etc/default/keyboard.bak",
    "${sudo} sed -i 's/XKBVARIANT=\"\w*"/XKBVARIANT=\"norman\"/g /etc/default/keyboard"]
    deps = ["#taskName", "^pkgName"]
    skip_if = ["which brew"] #only run if condition false
  
[task.norman__brew] # will run if brew is found instead
    cmds = ["@install norman"]
    deps = ["#taskName", "^pkgName"]
    run_if = ["which brew"] #only run if the result of run_if is true

[task.fish] # if both run_if and skip_if exist, then:
    run_if = ["which bash"]  #only run if condition true AND
    skip_if = ["which fish"] #run if condition false
```

## Installers

Default installers: 
```yaml
[installer.apt]
    run_if = ["which apt", "which apt-get"]
    sudo = true
    cmd =  "${sudo} apt install -y ${pkg}"

[installer.brew]
    run_if = ["which brew"]
    sudo = false
    cmd =  "${sudo} brew install ${pkg}"

[installer.apk]
    run_if = ["which apk"]
    sudo = false
    cmd =  "${sudo} apk add ${pkg}"

[installer.dnf]
    run_if = ["which dnf"]
    sudo = true
    cmd =  "${sudo} dnf install -y ${pkg}"

[installer.pacman]
    run_if = ["which pacman"]
    skip_if = ["which yay"]
    sudo = true
    cmd =  "${sudo} pacman -Syu ${pkg}"

[installer.yay]
    run_if = ["which yay"]
    sudo = true
    cmd =  "${sudo} yay -Syu ${pkg}"
```

To override an installer, include one of the above in your own configuration, and it will take precedence.

You can also create your own targets: 

```yaml
[installer.dnf]
    run_if = ["which dnf"]
    sudo = true # the default for this installer, over-ridden by command line args
    cmd =  "${sudo} dnf -y ${pkg}" #command line parameters

[installer.gvm] # golang version manager
    run_if = ["which gvm"]
    sudo = true # the default for this installer, over-ridden by command line args
    cmd =  "${sudo} gvm install ${pkg}" #command line parameters
```

And then all packages should be installable on that OS as specified.
//TODO: will have to supply a way to easily over-ride default package names with custom installer names

