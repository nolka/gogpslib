package test

import (
	"fmt"
	"testing"

	"github.com/nolka/gogpslib"
	"github.com/nolka/gogpslib/writer"
)

func TestGpxToPltWriteFile(t *testing.T) {
	gpx := gogpslib.GpxFormat{}
	plt := gogpslib.PltFormat{}
	gpx.Read("test_tracks/test1.gpx")

	s := writer.CreateStringWriter()

	plt.SetSegments(gpx.GetSegments())

	plt.WriteSegments(s)
	s.Write()

	fmt.Println(s.Content)
}
