# imagebeam

make images/gifs from a specific discord channel appear on an overlay for obs.

## usage

### prerequisites

make sure you have:

-   a discord bot token (with the bot joined in the desired server + appropriate permissions)
-   an obs browser source
-   golang >= 1.22
-   (optional) your's or friend's user discord ids, and/or specific channel ids for letting specific users/channels only being permitted

### installation / usage

#### step 1

open a terminal, and navigate into your local copy of the code:

```sh
cd /home/zeko/some_path/imagebeam
```

using `go`, compile the binary with

```sh
go build -o bin/ cmd/imagebeam-server.go
```

this will output a binary in the `bin/` directory, which should be something like `imagebeam-server` (or `imagebeam-server.exe`
if using windows)

#### step 2

in obs, modify the target browser source's url to be `http://127.0.0.1:8440`
you may also want to modify the width and height to be the same as your canvas's resolution, since scaling may distort or lower the quality of the displayed images.

![image of obs browser source settings](https://gist.github.com/assets/41507889/0eb8790d-9ea9-44a8-8557-f9a42af289c5)

#### step 3

with a terminal, you can then use imagebeam directly like so:

```sh
# run imagebeam without any specific permitted users/channels
bin/imagebeam-server -token <DISCORD_TOKEN>

# run with specific user(s) only allowed; multiple entries are split with a `,` (comma)
bin/imagebeam-server -token <DISCORD_TOKEN> -users 269609480212345678,1234567891234567890

# run with specific channels to be read from; can also have multiple entries, as well as combined with the `-users` argument
bin/imagebeam-server -token <DISCORD_TOKEN> -channels 280806440197867564
bin/imagebeam-server -token <DISCORD_TOKEN> -users 269609480212345678 -channels 280806440197867564,952337811013243546
```

it's probably not a good idea to run the binary with secrets in the command itself, so it may be more useful to use env variables through a setup script.

```sh
# maybe `./run.sh`
# this may depend on your shell of choice, but it can be something like this
set -l DISCORD_TOKEN '<discord token secret>'
set -l IB_USERS ''
set -l IB_CHANNELS ''

# then you can simply run the binary with the env variables instead
bin/imagebeam-server -token $DISCORD_TOKEN -users $IB_USERS -channels $IB_CHANNELS
```

#### step 4

run the server while obs is also open, and any images that the discord bot collects will be shown on the overlay.

##### note 1

wherever the bot is designated to read from, the bot supports displaying images/gifs from:

-   file attachments
-   linked images from discord's cdn (both cdn.discordapp.com and media.discordapp.net)
-   tenor links

##### note 2

the bot will only ever select the first attachment/link, so multiple links or attachments after the 1st will be ignored.

## project background

this was a small learning project to try and write in golang.
the idea came from [okayxairen in her discord](https://discord.com/channels/1142268839636258891/1142274158252785774/1215413075772055625), and since i was interesting in learning go, i thought this would be an ideal project to do.

some possible improvements after reflecting on the project:

-   the `net/http` web server is very likely underutilized
-   the discord<->overlay handling could be implemented more cleanly (maybe remove polling and go for some sort of event/messaging system)
-   a caching system for links/images already seen before
-   not implementing the overlay in a single html file
-   figuring out how to make proper compilation steps for cli programs, especially for cross-platform
-   continue learning how go projects are normally structured (i don't think this repo is that great...)

regardless, i still learned a fair bit and plan on implementing a handful of these on newer projects, plus more from along the way.
