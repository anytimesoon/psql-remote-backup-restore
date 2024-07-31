# Overview

Psqlbu is a command line util designed to connect to a remote postgres database, perform a backup and then restore the data to a local instance. 

> [!WARNING]  
> The restore function of this package drops the schema specified in the config file. This is by design so that the backup can be restored. Please take care when setting up the config.

# Getting started
Download the latest version from the releases page to a dedicated directory. 

Make the file executable with 
```
chmod 755 psqlbu-v0.x.0-[darwin/linux]-[amd64/arm64]
```

It is recommended to change the name of the executable to allow for easier use later. 
```
mv psqlbu-v0.x.0-[darwin/linux]-[amd64/arm64] psqlbu
```

#### Mac only
The executable is uncertified, so you will not be able to run it on mac without admin privileges.
To set the tool in quarantine and then adding the executable to your trusted apps, you will need to run the following
```
xattr -dr com.apple.quarantine psqlbu
spctl --add psqlbu
```

Add the directory to your $PATH to allow the command to be run from anywhere.

## Configuration
The util requires a configuration file to be located in the same directory as the executable, followed by the location where the command was run. By default, it will search for `main.yaml`, but a different file can be specified (see Flags)

An example config can be found [here](https://github.com/anytimesoon/psql-remote-backup-restore/blob/main/example/main_example.yaml)

## Flags

`-config [fileName]` - where file name is the name of the file without the `.yaml` extension. For example, to use the configuration in the file `custom_config.yaml` you should run 
```
psqlbu -config custom_config
```

## Args
There are two run time arguments which allow the user to bypass either the backup or restore portion of the operation. These will override the config.

`r` or `restore` - will only run the restore portion if the operation. It will look for the most recent backup file in the backups directory

`b` or `backup` - will only run the backup portion of the operation.
