[![Go Report Card](https://goreportcard.com/badge/github.com/jamesfcarter/galldir)](https://goreportcard.com/report/github.com/jamesfcarter/galldir)

# ![Galldir](http://jfc.org.uk/img/galldir_logo.png) galldir

The idea of this project is a photo gallery that is entirely driven from the
filesystem - no databases, no image upload interfaces, (almost) no
configuration.

Each directory within the filesystem is considered to be an album with all
image files within the directory being photos in the album. Albums may be
nested.

Each album has a human-readable name that is generated from the name of its
directory. This can be overridden by a text file called `.title` within the
directory.

Albums are displayed in date order, most recent first. The date of an album is
taken from the modification date of its directory, but this can be overridden
by a text file called `.date` within the directory that contains the album's
time in the format `YYY-MM-DD hh:mm:ss`. This is required when using S3 that
does not have directories per se.

![Galldir example album](http://jfc.org.uk/img/galldir_example.jpg)

## Installation

```
go get github.com/jamesfcarter/galldir/cmd/galldir
```

## Invocation

`galldir` requires two arguments: the address to listen on and the directory of
pictures to serve:
```
galldir -addr :3000 -dir ~/pictures
```

It is also possible to serve pictures from an S3 bucket:
```
galldir -addr :3000 -dir https://s3.eu-central-1.wasabisys.com/examplebucket
```

In both cases, browsing to http://localhost:3000/ would reach the gallery.

## License

This project is distributed under the [GNU GPL license
v3](https://www.gnu.org/licenses/gpl-3.0.en.html), see [LICENSE](./LICENSE) for
more information.

Uses [lightgallery.js](https://github.com/sachinchoolur/lightgallery.js) which
has its own
[license](https://github.com/sachinchoolur/lightgallery.js/blob/master/LICENSE.md).
