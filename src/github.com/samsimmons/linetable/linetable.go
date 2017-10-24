package main

import (
  "debug/elf"
  "debug/gosym"
  "flag"
  "log"
)

func main() {
  flag.Parse()
  path := flag.Arg(0)
  exe, err := elf.Open(path)
  if err != nil {
    log.Fatal(err)
  }
  var pclndat []byte
  if sec := exe.Section(".gopclntab"); sec != nil {
    pclndat, err = sec.Data()
    if err != nil {
      log.Fatalf("Cannot read .gopclntab section: %v", err)
    }
  }
  sec := exe.Section(".gosymtab")
  symTabRaw, err := sec.Data()
  pcln := gosym.NewLineTable(pclndat, exe.Section(".text").Addr)
  symTab, err := gosym.NewTable(symTabRaw, pcln)
  if err != nil {
    log.Fatal("Cannot create symbol table: %v", err)
  }
  sym := symTab.LookupFunc("main.main")
  filename, lineno, _ := symTab.PCToLine(sym.Entry)
  log.Printf("filename: %v\n", filename)
  log.Printf("lineno: %v\n", lineno)
}
