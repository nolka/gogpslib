package gogpslib

import (
	"time"

	"github.com/nolka/gogpslib/writer"
)

type FormatReaderWriter interface {
	Read(fileName string)
	GetSegments() []TrackSegment
	SetSegments(segments []TrackSegment)
	Write(w writer.Writer)
	GetGeoPointsCount() int
}

type GeoPoint struct {
	Lat        float64
	Lon        float64
	IsNewTrack bool
	Altitude   float64
	DateTime   time.Time

	Name        string
	Comment     string
	Description string
}

type TrackSegment struct {
	Name      string
	GeoPoints []GeoPoint
}

func (s *TrackSegment) GetGeoPoints() []GeoPoint {
	return s.GeoPoints
}

func (s *TrackSegment) GetPointsCount() int {
	return len(s.GeoPoints)
}
