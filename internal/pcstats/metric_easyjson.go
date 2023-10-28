package pcstats

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjsonEcc6f06fDecodeGithubComBazookajoe1MetricsCollectorInternalMetric(in *jlexer.Lexer, out *Metric) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "id":
			out.ID = string(in.String())
		case "type":
			out.MType = string(in.String())
		case "delta":
			if in.IsNull() {
				in.Skip()
				out.Delta = nil
			} else {
				if out.Delta == nil {
					out.Delta = new(int64)
				}
				*out.Delta = int64(in.Int64())
			}
		case "value":
			if in.IsNull() {
				in.Skip()
				out.Value = nil
			} else {
				if out.Value == nil {
					out.Value = new(float64)
				}
				*out.Value = float64(in.Float64())
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjsonEcc6f06fEncodeGithubComBazookajoe1MetricsCollectorInternalMetric(out *jwriter.Writer, in Metric) {
	out.RawByte('{')
	{
		const prefix string = ",\"id\":"
		out.RawString(prefix[1:])
		out.String(string(in.ID))
	}
	{
		const prefix string = ",\"type\":"
		out.RawString(prefix)
		out.String(string(in.MType))
	}
	if in.Delta != nil {
		const prefix string = ",\"delta\":"
		out.RawString(prefix)
		out.Int64(int64(*in.Delta))
	}
	if in.Value != nil {
		const prefix string = ",\"value\":"
		out.RawString(prefix)
		out.Float64(float64(*in.Value))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (m Metric) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjsonEcc6f06fEncodeGithubComBazookajoe1MetricsCollectorInternalMetric(&w, m)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (m Metric) MarshalEasyJSON(w *jwriter.Writer) {
	easyjsonEcc6f06fEncodeGithubComBazookajoe1MetricsCollectorInternalMetric(w, m)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (m *Metric) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjsonEcc6f06fDecodeGithubComBazookajoe1MetricsCollectorInternalMetric(&r, m)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (m *Metric) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjsonEcc6f06fDecodeGithubComBazookajoe1MetricsCollectorInternalMetric(l, m)
}
