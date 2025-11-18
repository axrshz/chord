package types

type Track struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ClientMessage struct {
	Action string `json:"action"` // "add" or "remove"
	Track  string `json:"track,omitempty"`
	ID     string `json:"id,omitempty"`
}

type ServerMessage struct {
	Type       string  `json:"type"` // "playlist"
	Playlist   []Track `json:"playlist"`
	NowPlaying int     `json:"nowPlaying"` // index of currently playing track (-1 = none)
}