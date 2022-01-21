Shoelace requires a configuration file to perform tasks. By default, it searches for a configuration file in:
* /usr/share/shoelace/default.toml
* $HOME/.config/shoelace/config.toml
* $XDG_HOME/shoelace/config.toml

## Usage

Shoelace can perform 3 different actions. Sync, install, and task. 
### Install

To perform an installation of a package through shoelace:
`shoelace install <pkgName>`

This will perform a lookup in the configuration file for package name substitutions, and then try to install the package.

### Sync
To perform a sync operation:
`shoelace sync <from> <to>`
This will symlink `from` into the `to` location. Optionally these values can be specified in a configuration file.

### Tasks
To perform a task operation:
`shoelace task <taskName>`
This will try run the specified task. The task needs to be defined in the configuration file loaded by shoelace.

#### Config File Simple Example
The simplest form is a single file with two sections:
```toml
[task.essential]
	installers = ["apt"]
	install = ["mercurial", "binutils", "bison", "build-essential"]
	
[installer.apt]
	sudo = false
	cmd =  "${sudo} apt install -y ${pkg}"
```

Then run shoelace.sh with `shoelace task essential`

## Tasks

Define what a task is here............

### Task Options

These options are listed in the order they are executed.

---
#### installers
Define which installers this task can be run with. If none of the installers are available, the task cannot run.
```toml
[task.example]
    installers = ["bash", "zsh"] 
```
---
#### run_if
Only run this task if the specified command returns true.
```toml
[task.example]
    run_if = ["which xcode"] 
```
---
#### skip_if
Skip this task if the command returns true.
```toml
[task.example]
    skip_if = ["which brew"] 
```
---
#### download
Download the specified file(s) from the internet to the target location(s)
```toml
[task.example]
    download = [["example.com/file.zip", "/home/example/file.zip"], ["example.com/file2.zip", "/tmp/file2.zip"]] 
```
---
#### deps
Install the required packages/tasks before running the install command. You can refer to other tasks here by prefixing the task name with a hash tag "#". Failure to install a dep, either as a task or as a package, prohibit this task from completing, and execution stops here.
```toml
[task.example]
    deps = ["git", "curl", "#custom_task"]
```
---
#### pre_cmd
Run the specified command before running the install command. This command can contain ASH and BASH variables, which will be substituted. If this command fails, execution is halted.
```toml
[task.example]
    pre_cmd = ["${CONFIG_PATH}/scripts/pre_install.sh"]
```
---
#### install
The package(s) to install.
```toml
[task.example]
    install = ["vim"]
```
---
#### post_cmd
Run the specified command after running the install command. This command can contain ASH and BASH variables, which will be substituted. If this fails, the installation is *not* rolled back.
```toml
[task.example]
    post_cmd = ["echo $USER", "echo HELLO WORLD!"]
```
---

## Shoelace variable substitution
Variables are available in the run_if, skip_if, download, pre_cmd, and post_cmd options.
* ORIGINAL_TASK  = Root task
* CURRENT_TASK   = Name of the currently executing task
* SUDO	       = If sudo should be enabled for that context
* CONFIG_PATH    = Full path location of the configuration file ? do we need paths for the various config files? packages.toml, ignores, etc?
* TARGET_PATH    = Target for symlinks
* SOURCE_PATH    = Source for symlinks

#### Available environment variables available in cmd lines
- sudo: if sudo should be enabled for commands
- pkg: pkg name
- installer: the name of the installer being used
- sudo: inserts sudo if enabled
- link_dest: the link destination for link creation
- link_src: the source directory containing original files to link to
- config_path: the path to the config file

## Packages
Package names, when defined in an `install` option for a task, are assumed to be the name used when installing that particular package. However, this can be overridden so that a single package name can resolve to platform/os specific package names.

In the following example, the `golang` package name is being defined as the following packages for each installer. shoelace will resolve which installer is being used, and if package name overrides exist for that installer, resolve the actual package name to install. Notice that the package name can also include version information in it as well, as long as the installer supports it.
```toml
[pkg.golang] 																			
    apt = "golang" 											
    apk = "google-go"  									
    yum = "go"
    yay = "golang@1.17"
    pacman = "go"
    gvm = "golang-1.17"
```

Package definitions can also contain a directive which defines either a single required installer, or an ordered list of installers. In the below case, we have defined that the only installer allowed for golang is gvm:
```toml
[pkg.golang] 												
    prefer = ["gvm"]
```

Or, we can define an ordered list of preferences:
```toml
[pkg.golang] 												
    prefer = ["gvm", "brew"]
```

## Installers
shoelace can support a multitude of various "installers", defined by a config. You can add your own installer just by adding a few lines. Below is an example of an installer with the required fields:

```toml
[installer.pacman]
	run_if = ["which pacman"]				
	cmd = "${sudo} pacman -S ${pkg}"
```

There are two variables required in a cmd line, namely `${sudo}` and `${pkg}`. This is further explained below.

### Installer Options

---
#### sudo
When using this installer, by default, run with sudo.
```toml
[installer.yay]
    sudo = true/false
```
---
#### run_if
Only use this installer if the detection condition is true.
```toml
[installer.yay]
   run_if = ["which yay"]
```
---
#### update
The command used by the installer to update it repo/cache information. This is run before the installer is used the first time.
```toml
[installer.yay]
   update = "${sudo} yay -Sy update"
```
---
#### cmd
The command to run when installing packages using this installer. Requires the `sudo` and `pkg` variables.
```toml
[installer.yay]
   cmd = "${sudo} yay -S ${pkg}"
```
---