# Jane CLI

This application is a very simple and slimmed down tool for using Jane via a CLI.  It
was created to help specifically to help practitioners who use screen readers to 
manager their own schedules as the Jane web app is not sufficiently setup to all
for screen reader integration.  This product is not sponsored by or associated with
Jane the company and exists because the company has failed to make adequate updates 
for accessibility.  This tool is rushed to solve the immediate needs of a single 
specific person; you are welcome to download and use the tool, but do so at your own
discretion.  I'm aiming to continue updating this over time, so feel free to create
issues for features you'd like to see, but know that this is currently a solo,
part-time development projects.

## Installation
Download the binary from the releases for your OS (only Windows and Linux are currently
supported) and install the executable along your PATH.  Copy the `config.yaml` file
to `etc/config.yaml` relative to where ever you put the executable and edit it to
include logs/debug and to point to your user yaml.

To launch the application, open up your favorite terminal and run `jane_cli`.  The first
time your launch the app, you'll need to do some setup with the `init` command.  All
subcommands provide some basic details via the `help` command.
