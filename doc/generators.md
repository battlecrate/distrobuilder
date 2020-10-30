# Generators

Generators are used to create, modify or remove files inside the rootfs.
Available generators are

* [cloud-init](#cloud-init)
* [dump](#dump)
* [hostname](#hostname)
* [hosts](#hosts)
* [remove](#remove)
* [template](#template)
* [upstart_tty](#upstart_tty)
* [lxd-agent](#lxd-agent)
* [fstab](#fstab)

In the image definition yaml, they are listed under `files`.

```yaml
files:
    - generator: <string> # which generator to use (required)
      name: <string>
      path: <string>
      content: <string>
      source: <string>
      template:
          properties: <map>
          when: <array>
      templated: <boolean>
      mode: <string>
      gid: <string>
      uid: <string>
      pongo: <boolean>
      architectures: <array> # filter
      releases: <array> # filter
      variants: <array> # filter
```

Filters can be applied to each entry in `files`.
Valid filters are `architecture`, `release` and `variant`.
See filters for more information.

## cloud-init

For LXC images, the generator disables cloud-init by disabling any cloud-init services, and creates the file `cloud-init.disable` which is checked by `cloud-init` on startup.

For LXD images, the generator creates templates depending on the provided name.
Valid names are `user-data`, `meta-data`, `vendor-data` and `network-config`.
The default `path` if not defined otherwise is `/var/lib/cloud/seed/nocloud-net/<name>`.
Setting `path`, `content` or `template.properties` will override the default values.

If `pongo` is true, the content will be processed using pongo2, and the context will be set appropriately (`{{ lxc.<variable> }}` or `{{ lxd.<variable> }}`).
See  [targets](targets.md).

## dump

The `dump` generator writes the provided `content` to a file set in `path`.
If provided, it will set the `mode` (octal format), `gid` (integer) and/or `uid` (integer).

If `pongo` is true, the content will be processed using pongo2, and the context will be set appropriately (`{{ lxc.<variable> }}` or `{{ lxd.<variable> }}`).
See  [targets](targets.md).

## hostname

For LXC images, the hostname generator writes the LXC specific string `LXC_NAME` to the hostname file set in `path`.
If the path doesn't exist, the generator does nothing.

For LXD images, the generator creates a template for `path`.
If the path doesn't exist, the generator does nothing.

## hosts

For LXC images, the generator adds the entry `127.0.0.1 LXC_NAME` to the hosts file set in `path`.

For LXD images, the generator creates a template for the hosts file set in `path`, adding an entry for `127.0.0.1 {{ container.name }}`.

## remove

The generator removes the file set in `path` from the container's root filesystem.

## template

This generator creates a custom LXD template.
The `name` field is used as the template's file name.
The `path` defines the target file in the container's root filesystem.
The `properties` key is a map of the template properties.

The `when` key can be one or more of:

* create (run at the time a new container is created from the image)
* copy (run when a container is created from an existing one)
* start (run every time the container is started)

If `pongo` is true, the content will be processed using pongo2, and the context will be set appropriately (`{{ lxc.<variable> }}` or `{{ lxd.<variable> }}`).
See  [targets](targets.md).

See [LXD image format](https://lxd.readthedocs.io/en/latest/image-handling/#image-format) for more.

## copy

The generator copies a file set in `source` to the container at `path`. All file specific properties such as `mode`, `gid` and `uid` are supported.

## upstart_tty

This generator creates an upstart job which prevents certain TTYs from starting.
The job script is written to `path`.

## lxd-agent

This generator creates the systemd unit files which are needed to start the lxd-agent in LXD VMs.

## fstab

This generator creates an /etc/fstab file which is used for VMs.
Its content is:

```
LABEL=rootfs  /         <fs>  <options>  0 0
LABEL=UEFI    /boot/efi vfat  defaults   0 0
```

The filesystem is taken from the LXD target (see [targets](targets.md)) which defaults to `ext4`.
The options are generated depending on the filesystem.
They cannot be overriden.
