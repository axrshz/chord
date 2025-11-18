package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/google/uuid"
	"github.com/olahol/melody"

	"github.com/axrshz/chord/types"
)

var (
	playlist   []types.Track // our in-memory playlist
	playlistMu sync.Mutex    // protects the playlist from concurrent writes
	m          *melody.Melody
)

func main() {
	m = melody.New()

	// Serve static files
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Home page
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Open /static/index.html to use the queue!")
	})

	// WebSocket endpoint
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		m.HandleRequest(w, r)
	})

	// When a client connects
	m.HandleConnect(func(s *melody.Session) {
		fmt.Println("Client connected")
		broadcastPlaylist() // send current state immediately
	})

	// When we receive a message
	m.HandleMessage(func(s *melody.Session, msg []byte) {
		var clientMsg types.ClientMessage
		if err := json.Unmarshal(msg, &clientMsg); err != nil {
			return // ignore malformed messages
		}

		playlistMu.Lock()
		switch clientMsg.Action {
		case "add":
			if clientMsg.Track != "" {
				newTrack := types.Track{
					ID:   uuid.New().String(),
					Name: clientMsg.Track,
				}
				playlist = append(playlist, newTrack)
				fmt.Printf("Added: %s\n", newTrack.Name)
			}
		case "remove":
			for i, t := range playlist {
				if t.ID == clientMsg.ID {
					playlist = append(playlist[:i], playlist[i+1:]...)
					fmt.Printf("Removed track ID %s\n", clientMsg.ID)
					break
				}
			}
		}
		playlistMu.Unlock()

		broadcastPlaylist()
	})

	fmt.Println("Server running → http://localhost:8080/static/index.html")
	http.ListenAndServe(":8080", nil)
}

// Send the full current playlist to ALL connected clients
func broadcastPlaylist() {
	playlistMu.Lock()
	defer playlistMu.Unlock()

	msg := types.ServerMessage{
		Type:       "playlist",
		Playlist:   playlist,
		NowPlaying: -1, // we'll use this later
	}

	data, _ := json.Marshal(msg)
	m.Broadcast(data)
}