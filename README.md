# GoBlogger

source code for my blog at [blog.rivanjuthani.com](https://blog.rivanjuthani.com/)

all your posts written in markdown go in the `/posts` folder, here is a template you can follow

```markdown
+++
title = "Just a demo"
description = "demo article written in markdown"
date = 2024-04-01
tags = ["example"]

[author]
name = "John Doe"
email = "john@example.com"
+++

# hello from markdown!
```

to generate your own certificates you can run this

```bash
openssl genrsa -out server.key 2048
openssl ecparam -genkey -name secp384r1 -out server.key
openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650
```

and put the two files in a folder called `/certs`

### Based on
- [Jon Calhoun's jonblog](https://github.com/joncalhoun/jonblog)