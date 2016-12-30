package main

import (
	"fmt"
	"os"

	"headerform"
)

func main() {
	doFixed()

	fmt.Println()

	doFlex()
}

func doFixed() {
	g := headerform.New(os.Stdout)
	g.InputDelimiter = ":"
	g.DividingLine = "-"

	// additional output possiblities:
	// g.OutputDelimiter = "|"
	// g.OutputDelimiter = " | "

	g.FormatAs("Username     |Passwd|   Uid| Gid| Name              |Home              |Shell      ")

	addRows(g)
	g.Print("long:x:70:70:a really long entry that is truncated:/sbin:/sbin/nologin")
}

func doFlex() {
	g := headerform.New(os.Stdout)
	g.InputDelimiter = ":"
	g.DividingLine = "-"

	// additional output possiblities:
	// g.OutputDelimiter = "|"
	// g.OutputDelimiter = " | "

	g.FormatAs("Username>|Passwd|<Uid|<Gid|<Name>|Home>|Shell      ")

	addRows(g)
	// flexible fields implicitly turn on buffering, so we need to explicitly emit
	// the header row and flush
	g.Headers() // calculates flexible columns, if any, prints headers
	g.Flush()   // flush the buffered rows

	g.Buffer = false // disable buffering, maintain previously calculated column widths

	// adding more rows use the previously calculated sizes
	g.Print("long:x:70:70:a really long flexible entry that is truncated:/sbin:/sbin/nologin")
}

func addRows(g *headerform.HeaderFormatter) {
	g.PrintStrings("root", "x", "0", "0", "root", "/root", "/bin/bash")

	g.PrintAny("bin", "x", 1, 1, "bin", "/bin", "/sbin/nologin")

	// uses InputDelimiter, with custom formatting for each field
	g.Printf("%s:%s:%03d:%03d:%s:%s:%s", "daemon", "y", 2, 2, "daemon", "/sbin", "/sbin/nologin")

	// uses InputDelimiter
	g.Print("adm:x:3:4:adm:/var/adm:/sbin/nologin")
	g.Print("lp:x:4:7:lp:/var/spool/lpd:/sbin/nologin")
	g.Print("sync:x:5:0:sync:/sbin:/bin/sync")
	g.Print("shutdown:x:6:0:shutdown:/sbin:/sbin/shutdown")
	g.Print("halt:x:7:0:halt:/sbin:/sbin/halt")
	g.Print("mail:x:8:12:mail:/var/spool/mail:/sbin/nologin")
	g.Print("uucp:x:10:14:uucp:/var/spool/uucp:/sbin/nologin")
	g.Print("operator:x:11:0:operator:/root:/sbin/nologin")
	g.Print("games:x:12:100:games:/usr/games:/sbin/nologin")
	g.Print("gopher:x:13:30:gopher:/var/gopher:/sbin/nologin")
	g.Print("ftp:x:14:50:FTP User:/var/ftp:/sbin/nologin")
	g.Print("nobody:x:99:99:Nobody:/:/sbin/nologin")
	g.Print("vagrant:x:65550:65550::/home/vagrant:/bin/bash")
	g.Print("saslauth:x:499:76:Saslauthd user:/var/empty/saslauth:/sbin/nologin")
	g.Print("vcsa:x:69:69:virtual console memory owner:/dev:/sbin/nologin")
	g.Print("postfix:x:89:89::/var/spool/postfix:/sbin/nologin")
	g.Print("sshd:x:74:74:Privilege-separated SSH:/var/empty/sshd:/sbin/nologin")
	g.Print("extra:x:25763:25761::/home/vagrant:/bin/bash:additional:fields")
}
