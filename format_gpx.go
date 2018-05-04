package gogpslib

import (
	"bytes"
	"encoding/xml"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/nolka/gogpslib/writer"
)

type link struct {
	Href string `xml:"href,attr"`
	Text string `xml:"text"`
	Type string `xml:"type"`
}

type author struct {
	Name string `xml:"name"`
	Link link   `xml:"link"`
}

type metadata struct {
	Name      string    `xml:"name"`
	Author    author    `xml:"author"`
	Copyright string    `xml:"copyright"`
	Link      link      `xml:"link"`
	Time      time.Time `xml:"time"`
}

type trkpt struct {
	Lat  string `xml:"lat,attr"`
	Lon  string `xml:"lon,attr"`
	Ele  string `xml:"ele"`
	Time string `xml:"time"`
}

type trkseg struct {
	Trkpt []trkpt `xml:"trkpt"`
}

type trk struct {
	Name   string `xml:"name"`
	Trkseg trkseg `xml:"trkseg"`
}

type gpx struct {
	Medatada metadata `xml:"metadata"`
	Trk      []trk    `xml:"trk"`
}

type GpxFormat struct {
	Gpx      gpx
	Segments []TrackSegment
}

func (s *GpxFormat) GetSegments() []TrackSegment {
	var segs []TrackSegment
	for _, _trk := range s.Gpx.Trk {
		var points []GeoPoint
		for _, pt := range _trk.Trkseg.Trkpt {
			point := GeoPoint{}
			point.Lat = StrToFloat(pt.Lat)
			point.Lon = StrToFloat(pt.Lon)

			if pt.Ele != "" {
				point.Altitude = StrToFloat(pt.Ele)
			}
			timeFormat := "2006-01-02T03:04:05Z"
			dt, err := time.Parse(timeFormat, pt.Time)
			if err == nil {
				point.DateTime = dt
			} else {
				log.Println(err)
			}
			points = append(points, point)
		}
		var seg TrackSegment
		seg.Name = _trk.Name
		seg.GeoPoints = points
		segs = append(segs, seg)
	}
	return segs
}

func (s *GpxFormat) SetSegments(segments []TrackSegment) {
	s.Segments = segments
}

func (s *GpxFormat) Read(fileName string) {
	file, err := os.Open(fileName)
	if err != nil {
		log.Print(err)
		return
	}
	defer file.Close()

	xmlBytes, _ := ioutil.ReadAll(file)

	err = xml.Unmarshal(xmlBytes, &s.Gpx)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *GpxFormat) WriteSegments(w writer.Writer) {
	tpl, err := template.New("gpxfile").Parse(GpxTemplate)
	if err != nil {
		log.Print(err)
		return
	}

	var b bytes.Buffer

	err = tpl.Execute(&b, s)
	if err != nil {
		log.Print(err)
		return
	}

	w.Append(b.String())
}

const GpxTemplate = `<?xml version="1.0" encoding="UTF-8"?>
<gpx
        version="1.0"
        creator="` + LibName + `"
        xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
        xmlns="http://www.topografix.com/GPX/1/0"
        xsi:schemaLocation="http://www.topografix.com/GPX/1/0 http://www.topografix.com/GPX/1/0/gpx.xsd">
	<metadata>
	</metadata>{{range .GetSegments}}
	<trk>
		<name>{{.Name}}</name>
		<trkseg>{{range .GeoPoints}}
			<trkpt lat="{{.Lat}}" lon="{{.Lon}}">
				{{if gt .Altitude 0.0}}<ele>{{.Altitude}}</ele>{{end}}
				<time>{{.DateTime.Format "2006-01-02T03:04:05Z"}}</time>
			</trkpt>{{end}}
		</trkseg>
	</trk>{{end}}
</gpx>`
