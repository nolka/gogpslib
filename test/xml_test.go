package test

import (
	"fmt"
	"testing"

	"github.com/nolka/gogpslib"
	"github.com/nolka/gogpslib/writer"
)

// func TestGpxReadFileTwo(t *testing.T) {
// 	const (
// 		SEG_COUNT = 2
// 		POINT_COUNT1 = 3
// 		POINT_COUNT2 = 4
// 	)
//
// 	gpx := converter.GpxFormat{}
// 	gpx.Read("test_tracks/test1.gpx")
//
// 	s := gpx.GetSegments()
// 	j, _ := json.Marshal(&s)
// 	fmt.Print(string(j))

// if len(segments) != SEG_COUNT {
// 	t.Fatal("Segments count must be 2", fmt.Sprintf("Actual: [%d]", len(segments)))
// }
//
// segmentOne := segments[0]

// if len(segmentOne.GeoPoints) != POINT_COUNT1 {
// 	t.Fatal("Points count does not match in segment 1", fmt.Sprintf("Actual: [%d]", len(segmentOne.GeoPoints)))
// }
//
// if !segmentOne.GeoPoints[0].IsNewTrack {
// 	t.Fatal("First point in track must be a track break")
// }
//
// if segmentOne.GeoPoints[1].IsNewTrack == true {
// 	t.Fatal("First point in track must NOT be a track break")
// }
//
// segmentTwo := segments[1]
//
// if len(segmentTwo.GeoPoints) != POINT_COUNT2 {
// 	t.Fatal("Points count does not match in segment 2", fmt.Sprintf("Actual: [%d]", len(segmentTwo.GeoPoints)))
// }
//
// if !segmentTwo.GeoPoints[0].IsNewTrack {
// 	t.Fatal("First point of secont segment, Fourth point in track must be a track break")
// }
//
// if segmentTwo.GeoPoints[0].Lat != 55.809181 {
// 	t.Fatal("Segment two: Incorrect Latitude parsed. Actual is 55.809181")
// }
// if segmentTwo.GeoPoints[0].Lon != 92.505028 {
// 	t.Fatal("Segment two: Incorrect Latitude parsed. Actual is 92.505028")
// }
// }

func TestGpxWriteFile(t *testing.T) {
	gpx := gogpslib.GpxFormat{}
	gpx.Read("test_tracks/test1.gpx")

	// j, _ := json.Marshal(&gpx)
	// fmt.Print(string(j))

	s := writer.CreateStringWriter()

	gpx.WriteSegments(s)
	s.Write()

	fmt.Println(s.Content)
}
