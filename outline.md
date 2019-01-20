# What the dep?

## Prehistory

- Not going to dwell too long on this, but it's no bad thing to understand how
  we got here
- `$GOPATH`, love it or hate it
- Zeroth generation: `go get` was the only mechanism
    - Challenging for anyone writing code in a 2010s style with lots of
      dependencies
    - Could only have `master`, without really interesting domain fronting
      hackery like `gopkg.in`

## The clone wars

- The first generation of tooling: Glide, godep, rubigo, et al
- Those with ambitions of being the standard package manager generally mimicked
  other languages: Glide, for example, explicitly called out Cargo, npm,
  Composer, and a few others as direct inspirations and comparisons
- `$GO15VENDOREXPERIMENT` helped a lot here: in 2015, you could vendor your
  dependencies into `vendor` directories, which had their own special rules for
  import resolution, and this provided an easier way for those tools to work
- Nevertheless, there were definite challenges on a tooling level: while most
  languages with widely used package managers actually decouple manager
  development from the language (with exceptions like Rust), Go didn't even
  really provide some of the primitives required; you were always hacking
  around limitations in the tooling, because of how `$GOPATH` centric it was
- The profileration of alternatives also made it hard to coalesce around one
  ecosystem: each option had its own metadata format, so if you imported a
  package that didn't use your system, then you had to cross your fingers that
  whatever inference it did around branch and tag names worked reliably
- Finally, there were quality issues: Glide probably should have won, in my
  opinion, but for glaring issues around updating existing dependencies

## `dep`

- For years, the biggest issue named in the annual Go developer survey was
  package and dependency management
- Sam Boyer went away and did a shitload of research and came up with a design
  for what became `dep`
- Interestingly, even with `dep` in a fairly rough beta state, when the 2017
  survey was performed, generics had become the #1 issue, presumably because
  Go developers generally considered that package management was about to be
  "solved"
- `dep init` included a really interesting feature called "inference", where it
  would work really hard to figure out from the code what versions of what
  dependencies you had. It was kind of magic
- `dep` still relied on `$GOPATH`, ultimately

## `vgo`

- While `dep` was out getting community consensus, Russ Cox on the Go core team
  thought he could do better
- Enter `vgo`!
- `vgo` took a different tack: rather than being a standalone tool that was
  blessed by the Go team (and shipped with Go), `vgo` would be deeply
  integrated _into_ the `go` command line tool
- Unlike `dep`, `vgo` also mostly got rid of the `$GOPATH` construct from the
  perspective of the using developer: instead, you end up with a developer
  experience that's closer to other languages whereby the project root is
  effectively defined as "the place where the package metadata lives"
- Most package managers use SAT solvers: figuring out the full dependency tree
  is NP-complete, and you end up with the slow behaviour you get in most such
  tools
- `vgo` uses what Russ Cox calls MVS: minimal version selection. It's an
  attempt to cut down the problem space by only installing the actual versions
  in the dependency metadata wherever possible

## Today

- For better or worse, Go 1.11 (released in August) shipped `vgo` as an
  integral (although nominally experimental) part of the `go` command line
  tool, effectively killing `dep` (since many of the assumptions `dep` had
  worked under no longer hold true).  You might as well get on board, no matter
  your thoughts on the whole drama surrounding `dep` and `vgo` (which I mostly
  didn't even get into)
- You opt into the new world by having a `go.mod` file at your module root.
  Generally this is created via `go mod init`
  - Define a module vs a package
- Using `go mod` _requires_ semantic versioning, both of your code and its
  dependencies. (The Go ecosystem is new enough that this is mostly the case
  anyway, but it's something to be aware of.)
- `go.mod` files use their own special syntax (of course), but it's human
  readable and, in a pinch, writable (don't)
- No lock file because of the use of MVS: the same `go.mod` should always
  provide the same set of package versions, barring shenanigans on the part of
  a module author with redefining tags or force-pushing branches
- There _is_ a `go.sum` file used for crypto purposes to detect shenanigans;
  you should commit it, but it's not a lock file per se

## Examples

### Existing project

- `go mod init` performs a bunch of magic to try to migrate existing dependency
  metadata (including `dep`, Glide, and a bunch of others) and to use your Git
  remote set up to try to intuit your module path
  - Use gpsd as an example
- Set up an example with no metadata

### New clone

- Because it's Go, there's no real separate fetch step: just `go build` or
  `go install` and it'll fetch whatever `go.mod` tells it to

### Upgrading dependencies

- `go get -u` now takes version numbers with `@x.y.z` suffixes

### Releasing a project

- Must be semver
- Must use the `vgo` tagging rules: versions must be tagged with `vX.Y.Z` and
  not `X.Y.Z`, which is mildly annoying if you were doing it the other way
