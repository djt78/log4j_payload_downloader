package main

import (
  "fmt"
  "os"
  "time"
  "io/ioutil"
  "crypto/md5"
  "net/url"
  "net/http"
  "strings"
  "github.com/go-ldap/ldap"
)

func use(vals ...interface{}) {
    for _, val := range vals {
        _ = val
    }
}

func grab(entry *ldap.Entry) {
   addr, err := url.Parse(entry.GetAttributeValue("javaCodeBase"))
   if err != nil {
	   fmt.Println("Unable to get class url")
	   return
   }

   x := fmt.Sprintf("/%s.class",entry.GetAttributeValue("javaFactory"))
   addr.Path = x
   fmt.Println("Class URL:"+addr.String());

   client := http.Client {
	   Timeout: 10 * time.Second,
   }

   response,err := client.Get(addr.String())
   defer response.Body.Close()

   bytes, err := ioutil.ReadAll(response.Body)
   if err != nil {
     fmt.Println("Unable to get bytes from response body")
     return
   }

   savedir :="downloaded/"
   err = os.MkdirAll(savedir,os.ModePerm)
   if err != nil {
	   savedir = ""
   }

   MD5 := fmt.Sprintf("payload_%x.class",md5.Sum(bytes))
   
   fmt.Println("Saving Class: "+MD5)
   of, err := os.Create(savedir + MD5)
   
   of.Write(bytes)
   of.Close()
   
}

func DownloadFromLdap(addr string) {

  u, err := url.Parse(addr)
  if err!= nil {
    fmt.Println("Cannot parse URL:"+addr);
  }

  fmt.Println("Scheme:"+u.Scheme)
  fmt.Println("Host:"+u.Host)
  fmt.Println("Path:"+u.Path)
  fmt.Println("")

  lc, err := ldap.DialURL(u.String())
  if err != nil {
     fmt.Println("Cannot Connect to:"+addr)
     return
  }

  defer lc.Close()

  err = lc.UnauthenticatedBind("")
  if err != nil {
    fmt.Println("Cannot LDAP bind to:"+addr)
    return;
  }


  sreq := ldap.NewSearchRequest(
	  strings.TrimLeft(u.Path,"/"), 
	  ldap.ScopeBaseObject,
	  ldap.DerefAlways,
	  0,
	  0,
	  false,
	  "(ObjectClass=*)",
	  []string{},
	  []ldap.Control {
             &ldap.ControlManageDsaIT{},
	  },
  )


  sres, err := lc.Search(sreq)

  if err != nil {
    fmt.Println("Search Failed")
    return
  }

  //files := []string{}

  for key, value := range sres.Entries {
    fmt.Println("Key[",key,"] Value[",value,"]")
    class := value.GetAttributeValue("objectClass")

    if class == "javaNamingReference" {
	    fmt.Println("objectClass: "+class)
	    grab(value);
    } else {
	    fmt.Println("Not JavaNamingReference... "+class)
    }

  }




  use(sres)

  fmt.Println("")
  fmt.Println("")
  fmt.Println("")
  fmt.Println("")


  lc.Close()

}






func main() {

  argc := len(os.Args)
  if (argc != 2) {
	  fmt.Println("Usage: l4jdl URL\n");
	  fmt.Println("Example: l4jdl ldap://somedomain/Exploit\n");
	  os.Exit(3)
  }
  url := os.Args[1]   

  DownloadFromLdap(url)

}
