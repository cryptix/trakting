package views

import (
	"fmt"

	"github.com/cryptix/trakting/types"
	"github.com/soroushjp/humble"
)

type Track struct {
	humble.Identifier

	Track  types.Track
	Parent *TrackList
}

func (p *Track) RenderHTML() string {
	return fmt.Sprintf(`<li>%s</li>
		`, p.Track.Name)
}

func (p *Track) OuterTag() string {
	return "div"
}
