package main

import "fmt"
//import "io"
import "github.com/amarburg/lazyfs"
import "github.com/amarburg/go-quicktime"

var TestUrlRoot = "http://localhost:8080/files/"
var TestMovPath = "CamHD_Vent_Short.mov"

var SparseHttpStoreRoot = "test_files/httpsparse/"

func main() {

  source,err := lazyfs.OpenHttpFSSource(TestUrlRoot)
  if err != nil {
    panic("Couldn't open HttpFSSource")
  }

  store,err := lazyfs.OpenSparseFileFSStore( SparseHttpStoreRoot )
  if store == nil {
    panic("Couldn't open SparesFileFSStore")
  }

  source.SetBackingStore( store )

  file,err := source.Open( TestMovPath )
  if err != nil {
    panic("Couldn't open AlphabetPath")
  }

  // var offset int64 = 0
  // var indent = 0

  sz,_ := file.FileSize()
//  ParseAtom( file, offset, sz, indent )

  tree := quicktime.BuildTree( file, sz )

  quicktime.DumpTree( file, tree )

  moov := tree.FindAtom("moov")
  if moov == nil { panic("Can't find MOOV atom")}

  tracks := moov.FindAtoms("trak")
  if tracks == nil || len(tracks) == 0 { panic("Couldn't find any TRAKs in the MOOV")}
  fmt.Println("Found",len(tracks),"TRAK atoms")

  var track *quicktime.Atom = nil
  for i,t := range tracks {
    fmt.Println(t, t.Type)
    mdia := t.FindAtom("mdia")
    if mdia == nil {
      fmt.Println("No mdia track",i)
      continue
    }

    minf := mdia.FindAtom("minf")
    if minf == nil {
      fmt.Println("No minf track",i)
      continue
    }

    if minf.FindAtom("vmhd") != nil {
      fmt.Println("Found vmhd")
      track = t
      break
    }
  }

  if track == nil { panic("Couldn't identify the Video track")}

  stbl_atom := track.FindAtom("mdia").FindAtom("minf").FindAtom("stbl")
  stbl,_ := quicktime.ParseSTBL( stbl_atom )

  fmt.Println("Found track with video information")

  // Find movie length
  num_frames := stbl.NumFrames()

  fmt.Println("Movie has",num_frames,"frames")

  fmt.Println("Chunk table:")
  for idx,offset := range stbl.Stco.ChunkOffsets {
    fmt.Printf("   %d %20d\n",idx+1, offset )
  }

  //fmt.Println(stbl)

  for sample := 1; sample <= num_frames; sample++ {
      chunk,chunk_start,relasample := stbl.Stsc.SampleChunk( sample )
      fmt.Println("Sample", sample,"is in chunk",chunk,"the",relasample,"'th sample; the chunk starts at sample",chunk_start)

      offset,_ := stbl.SampleOffset( sample )
      fmt.Println("Sample at byte",offset,"in file")


  }


  // Try extracting a frame
  frame := 1
  LoadFrame( frame, stbl, file )
}


func LoadFrame( frame int, stbl quicktime.STBLAtom, file io.ReaderAt ) {

  frame_offset,frame_size,_ := stbl.SampleOffsetSize( frame )

  fmt.Printf("Extracting frame %d at offset %d size %d\n", frame, frame_offset, frame_size)

  buf := make([]byte, frame_size)
  n,err = file.ReadAt( buf, frame_offset )

  if n != frame_size { panic(fmt.Sprintf("Tried to read %d bytes but got %d instead",frame_size,n))}

  
}