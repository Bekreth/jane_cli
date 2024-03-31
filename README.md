# Jane CLI

This application is a very simple and slimmed down tool for using Jane via a CLI.  It
was created to help specifically to help practitioners who uses a screen reader to 
manager their own schedules as the Jane web app is not sufficiently setup to all
for screen reader integration.  This product is not sponsored by or associated with
Jane the company and exists because the company has failed to make adequate updates 
for accessibility.  This tool is rushed to solve the immediate needs of a single 
specific person; you are welcome to download and use the tool, but do so at your own
discretion.  I'm aiming to continue updating this over time, so feel free to create
issues for features you'd like to see, but know that this is currently a solo,
part-time development projects.

# Installation
## Windows
Download the MSI for the version you're using and launch that.

## Unix People
Download the binary from the releases for your OS (only Windows and Linux are currently
supported) and install the executable along your PATH.  Copy the `config.yaml` file
to `etc/config.yaml` relative to where ever you put the executable and edit it to
include logs/debug and to point to your user yaml.

# How to Use
To launch the application, open up your favorite terminal and run `jane_cli`.  The first
time your launch the app, you'll need to do some setup with the `init` command.  All
subcommands provide some basic details via the `help` command.

## Warning
While this application is functional with screen readers, its not without some minor
issues. Moving your cursor forwards and backwards through the text reads out the letters
the wrong way round, and there's no way yet to skip forwards/backwards full words and
have them read.  This tool is meant as a stop gap and will allow you to use Jane
more easily than their web portal, but it has some weaknesses.

# Features
* See schedule for specific days
* Book and cancel appointments
* Read and write Chart Notes
