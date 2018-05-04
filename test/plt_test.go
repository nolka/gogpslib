package test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/nolka/gogpslib"
	"github.com/nolka/gogpslib/writer"
)

func TestReadFileOne(t *testing.T) {
	plt := gogpslib.PltFormat{}
	plt.Read("test_tracks/test1.plt")
	segments := plt.GetSegments()

	points := segments[0].GeoPoints

	if len(segments) != 1 {
		t.Fatal("Incorrect segments count", fmt.Sprintf("Actual: [%d]", len(segments)))
	}

	if plt.FormatHeaders.Version != "OziExplorer Track Point File Version 2.1" {
		t.Fatal("Plt version parsed incorrect", fmt.Sprintf("Actual: [%s]", plt.FormatHeaders.Version))
	}

	if plt.FormatHeaders.GeodeticDatum != "WGS 84" {
		t.Fatal("Geodetic Datum parsed incorrect")
	}

	if plt.FormatHeaders.Reminder != "Altitude is in Feet" {
		t.Fatal("Reminder parsed incorrect")
	}

	if plt.FormatHeaders.Reserved != "Reserved 3" {
		t.Fatal("Reserved parsed incorrect")
	}

	if plt.PointsCount != 3 {
		t.Fatal(errors.New("Points count parsed incorrect"))
	}

	if len(points) != plt.PointsCount {
		t.Fatal("Points count does not match", fmt.Sprintf("Actual: [%d]", len(points)))
	}

	if !points[0].IsNewTrack {
		t.Fatal("First point in track must be a track break")
	}

	if points[1].IsNewTrack == true {
		t.Fatal("First point in track must NOT be a track break")
	}
}

func TestReadFileTwo(t *testing.T) {
	const (
		SEG_COUNT    = 2
		POINT_COUNT1 = 3
		POINT_COUNT2 = 4
	)

	plt := gogpslib.PltFormat{}
	plt.Read("test_tracks/test2.plt")
	segments := plt.GetSegments()

	if len(segments) != SEG_COUNT {
		t.Fatal("Segments count must be 2", fmt.Sprintf("Actual: [%d]", len(segments)))
	}

	segmentOne := segments[0]

	if len(segmentOne.GeoPoints) != POINT_COUNT1 {
		t.Fatal("Points count does not match in segment 1", fmt.Sprintf("Actual: [%d]", len(segmentOne.GeoPoints)))
	}

	if !segmentOne.GeoPoints[0].IsNewTrack {
		t.Fatal("First point in track must be a track break")
	}

	if segmentOne.GeoPoints[1].IsNewTrack == true {
		t.Fatal("First point in track must NOT be a track break")
	}

	segmentTwo := segments[1]

	if len(segmentTwo.GeoPoints) != POINT_COUNT2 {
		t.Fatal("Points count does not match in segment 2", fmt.Sprintf("Actual: [%d]", len(segmentTwo.GeoPoints)))
	}

	if !segmentTwo.GeoPoints[0].IsNewTrack {
		t.Fatal("First point of secont segment, Fourth point in track must be a track break")
	}

	if segmentTwo.GeoPoints[0].Lat != 55.809181 {
		t.Fatal("Segment two: Incorrect Latitude parsed. Actual is 55.809181")
	}
	if segmentTwo.GeoPoints[0].Lon != 92.505028 {
		t.Fatal("Segment two: Incorrect Latitude parsed. Actual is 92.505028")
	}
}

func TestWriteFile(t *testing.T) {
	plt := gogpslib.PltFormat{}
	plt.Read("test_tracks/test2.plt")

	s := writer.CreateStringWriter()
	f := writer.CreateFileWriter("/tmp/plt")

	plt.WriteSegments(s)
	s.Write()

	plt.WriteSegments(f)
	f.Write()

	fmt.Println(s.Content)
}
