package main

import (
  "fmt"
  "os"
  "io"
  "os/exec"
  "strings"
  "archive/tar"
  "path/filepath"
)

func main() {
  if len(os.Args) <= 2 {
    panic("[binary] [output.tar]")
  }
  file := os.Args[1]
  outfile := os.Args[2]

  names, err := linkedFiles(file)
  if err != nil {
    panic(err)
  }

  fi, err := os.Create(outfile)
  if err != nil {
    panic(err)
  }
  defer fi.Close()

  tw := tar.NewWriter(fi)
  defer tw.Close()

  addFile(tw, file)


  for _, name := range names {
    addFile(tw, name)
  }
}

func addFile(tw *tar.Writer, name string) error {
  fmt.Println(name)
  // open files for taring
  path, err := filepath.EvalSymlinks(name)
  if err != nil {
    return err
  }

  f, err := os.Open(path)
  if err != nil {
    return err
  }
  defer f.Close()

  fi, err := f.Stat()
  if err != nil {
    return err
  }

  fm := fi.Mode()
  header := &tar.Header{
    Typeflag: tar.TypeReg,
    Name: strings.TrimPrefix(name, "/"),
    ModTime: fi.ModTime(),
    Mode: int64(fm.Perm()), // or'd with c_IS* constants later
    Size: fi.Size(),
  }
  if err := tw.WriteHeader(header); err != nil {
    return err
  }
  // copy file data into tar writer
  if _, err := io.Copy(tw, f); err != nil {
    return err
  }
  return nil
}

func linkedFiles(file string) ([]string, error) {
  var names []string
  o, err := exec.Command("ldd", file).Output()
  if err != nil {
    return nil, err
  }
  for _, p := range strings.Split(string(o), "\n") {
    f := strings.Split(p, " ")
    if len(f) >= 3 && f[1] == "=>" && len(f[2]) > 0 {
      names = append(names, f[2])
    } else if len(f) > 0 {
        // ld
      path := strings.TrimSpace(f[0])
      if _, err := os.Stat(path); err == nil {
        names = append(names, path)
      }
    }
  }
  return names, nil
}
