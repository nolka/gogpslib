package gogpslib

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/nolka/gogpslib/writer"
)

type TrackDetails struct {
	TrackWidth     int
	TrackColor     int
	Description    string
	TrackSkip      int
	TrackType      int
	TrackFillStyle int
	TrackFillColor int
}

type FormatHeaders struct {
	Version       string
	GeodeticDatum string
	Reminder      string
	Reserved      string
}

type PltFormat struct {
	FormatHeaders *FormatHeaders
	Details       *TrackDetails
	PointsCount   int
	Segments      []TrackSegment
}

func (s *PltFormat) GetSegments() []TrackSegment {
	return s.Segments
}

func (s *PltFormat) SetSegments(segments []TrackSegment) {
	s.Segments = segments
}

func (s *PltFormat) GetGeoPointsCount() int {
	var count int = 0
	for _, s := range s.GetSegments() {
		count += s.GetPointsCount()
	}
	return count
}

func (s *PltFormat) GetDefaultDetails() *TrackDetails {
	return &TrackDetails{
		1,
		255,
		LibName,
		1,
		1,
		0,
		0,
	}
}

func (s *PltFormat) GetDefaultHeaders() *FormatHeaders {
	return &FormatHeaders{
		"OziExplorer Track Point File Version 2.1",
		"WGS 84",
		"Altitude is in Feet",
		"Reserved 3",
	}
}

func (s *PltFormat) ParseDetails(detailsLine string) *TrackDetails {
	parts := strings.Split(detailsLine, ",")
	details := TrackDetails{}
	details.TrackWidth = StrToInt(parts[1])
	details.TrackColor = StrToInt(parts[2])
	details.Description = parts[3]
	details.TrackSkip = StrToInt(parts[4])
	details.TrackType = StrToInt(parts[5])
	details.TrackFillStyle = StrToInt(parts[5])
	details.TrackFillColor = StrToInt(parts[6])
	return &details
}

func (s *PltFormat) Read(fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		log.Print(err)
		return
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	s.FormatHeaders = s.GetDefaultHeaders()
	s.FormatHeaders.Version, _ = reader.ReadString('\n')
	s.FormatHeaders.GeodeticDatum, _ = reader.ReadString('\n')
	s.FormatHeaders.Reminder, _ = reader.ReadString('\n')
	s.FormatHeaders.Reserved, _ = reader.ReadString('\n')

	s.FormatHeaders.Version = strings.Trim(s.FormatHeaders.Version, "\r\n")
	s.FormatHeaders.GeodeticDatum = strings.Trim(s.FormatHeaders.GeodeticDatum, "\r\n")
	s.FormatHeaders.Reminder = strings.Trim(s.FormatHeaders.Reminder, "\r\n")
	s.FormatHeaders.Reserved = strings.Trim(s.FormatHeaders.Reserved, "\r\n")

	detailsLine, _ := reader.ReadString('\n')
	s.Details = s.ParseDetails(detailsLine)

	pointsCountLine, _ := reader.ReadString('\n')
	s.PointsCount = StrToInt(pointsCountLine)

	scanner := bufio.NewScanner(reader)

	var segment TrackSegment
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), ",")

		pointDate := time.Now()
		if parts[4] != "" {
			pointDate = ParseDelphiDate(parts[4])
		}

		var point GeoPoint = GeoPoint{
			StrToFloat(parts[0]),
			StrToFloat(parts[1]),
			func() bool {
				if parts[2] == "1" {
					return true
				}
				return false
			}(),
			StrToFloat(parts[3]),
			pointDate,
			"",
			"",
			"",
		}

		if point.IsNewTrack {
			if len(segment.GeoPoints) > 0 {
				s.Segments = append(s.Segments, segment)
			}
			segment = TrackSegment{}
		}
		segment.GeoPoints = append(segment.GeoPoints, point)
	}
	if len(segment.GeoPoints) > 0 {
		s.Segments = append(s.Segments, segment)
	}
}

func (s *PltFormat) WriteSegments(w writer.Writer) {
	if s.FormatHeaders == nil {
		s.FormatHeaders = s.GetDefaultHeaders()
	}
	w.Append(s.FormatHeaders.Version + "\r\n")
	w.Append(s.FormatHeaders.GeodeticDatum + "\r\n")
	w.Append(s.FormatHeaders.Reminder + "\r\n")
	w.Append(s.FormatHeaders.Reserved + "\r\n")

	if s.Details == nil {
		s.Details = s.GetDefaultDetails()
	}
	d := s.Details

	w.Append(fmt.Sprintf("0,%d, %d,%s,%d,%d,%d, %d\r\n", d.TrackWidth, d.TrackColor, d.Description, d.TrackSkip, d.TrackType, d.TrackFillStyle, d.TrackFillColor))
	w.Append(fmt.Sprintf(" %d\r\n", s.GetGeoPointsCount()))
	for _, seg := range s.GetSegments() {
		for i, p := range seg.GetGeoPoints() {
			var isNewTrack int = 0
			if p.IsNewTrack || i == 0 {
				isNewTrack = 1
			}
			if p.Altitude == 0 {
				p.Altitude = -777
			}
			w.Append(fmt.Sprintf("%f,%f,%d,%f,%f,,\r\n", p.Lat, p.Lon, isNewTrack, p.Altitude, ToDelphiDate(p.DateTime)))
		}
	}
}
