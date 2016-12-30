# headerform text table layout in go

Descriptive text table formatting for command line/terminal go programs.

Describe the layout using a literal, visual format, then feed strings into it.

Both fixed width (explicitly defined by the layout) and flexible width (sized base on the maximum
length of the data in a column) are supported.  A single layout can use both fixed and flexible columns
in the same output.

See `src/main.go` for sample code.

See the documentation for the `FormatAs` function for more detail on the layout specification.

## TODO

This repo layout is lame and is stupid for use with `go get`.

## Fixed Layout

Fixed layout is done by describing the layout with aligned headers and delimiters, where whitespace defines alignment and padding:
```
"Username     |Passwd|   Uid| Gid| Name              |Home              |Shell      "
```
That layout produces the following output:
```
Username      Passwd    Uid  Gid        Name         Home               Shell      
------------- ------ ------ ---- ------------------- ------------------ -----------
root          x           0    0        root         /root              /bin/bash  
bin           x           1    1         bin         /bin               /sbin/nologin
daemon        y         002  002       daemon        /sbin              /sbin/nologin
adm           x           3    4         adm         /var/adm           /sbin/nologin
lp            x           4    7         lp          /var/spool/lpd     /sbin/nologin
sync          x           5    0        sync         /sbin              /bin/sync  
shutdown      x           6    0      shutdown       /sbin              /sbin/shutdown
halt          x           7    0        halt         /sbin              /sbin/halt 
mail          x           8   12        mail         /var/spool/mail    /sbin/nologin
uucp          x          10   14        uucp         /var/spool/uucp    /sbin/nologin
operator      x          11    0      operator       /root              /sbin/nologin
games         x          12  100        games        /usr/games         /sbin/nologin
gopher        x          13   30       gopher        /var/gopher        /sbin/nologin
ftp           x          14   50      FTP User       /var/ftp           /sbin/nologin
nobody        x          99   99       Nobody        /                  /sbin/nologin
vagrant       x       65550 …550                     /home/vagrant      /bin/bash  
saslauth      x         499   76   Saslauthd user    /var/empty/saslau… /sbin/nologin
vcsa          x          69   69 virtual console me… /dev               /sbin/nologin
postfix       x          89   89                     /var/spool/postfix /sbin/nologin
sshd          x          74   74 Privilege-separate… /var/empty/sshd    /sbin/nologin
extra         x       25763 …761                     /home/vagrant      /bin/bash   additional fields
long          x          70   70 a really long entr… /sbin              /sbin/nologin
```

## Flexible Layout

Flexible layout is specified using `<` and `>` around the fields that should resize based on content.
```
"Username>|Passwd|<Uid|<Gid|<Name>|Home>|Shell      "
```
Which produces output like the following:
```
Username Passwd   Uid   Gid             Name             Home                Shell      
-------- ------ ----- ----- ---------------------------- ------------------- -----------
root     x          0     0             root             /root               /bin/bash  
bin      x          1     1             bin              /bin                /sbin/nologin
daemon   y        002   002            daemon            /sbin               /sbin/nologin
adm      x          3     4             adm              /var/adm            /sbin/nologin
lp       x          4     7              lp              /var/spool/lpd      /sbin/nologin
sync     x          5     0             sync             /sbin               /bin/sync  
shutdown x          6     0           shutdown           /sbin               /sbin/shutdown
halt     x          7     0             halt             /sbin               /sbin/halt 
mail     x          8    12             mail             /var/spool/mail     /sbin/nologin
uucp     x         10    14             uucp             /var/spool/uucp     /sbin/nologin
operator x         11     0           operator           /root               /sbin/nologin
games    x         12   100            games             /usr/games          /sbin/nologin
gopher   x         13    30            gopher            /var/gopher         /sbin/nologin
ftp      x         14    50           FTP User           /var/ftp            /sbin/nologin
nobody   x         99    99            Nobody            /                   /sbin/nologin
vagrant  x      65550 65550                              /home/vagrant       /bin/bash  
saslauth x        499    76        Saslauthd user        /var/empty/saslauth /sbin/nologin
vcsa     x         69    69 virtual console memory owner /dev                /sbin/nologin
postfix  x         89    89                              /var/spool/postfix  /sbin/nologin
sshd     x         74    74   Privilege-separated SSH    /var/empty/sshd     /sbin/nologin
extra    x      25763 25761                              /home/vagrant       /bin/bash   additional fields
long     x         70    70 a really long flexible entr… /sbin               /sbin/nologin
```

In this case, some set of rows are feed in and the output is flushed (and the widths calculated) and then more rows are emitted.
The name column in the last row above is truncated because an explicit flush was used to fix the column widths to the data seen
so far.



