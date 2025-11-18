package main

import (
	"fmt"
	"net/http"

	"github.com/olahol/melody" // ← new import
)

func main() {
	// Create a new Melody instance (this is our WebSocket manager)
	m := melody.New()

	// ---------- Regular HTTP route (still works!) ----------
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, synced music queue! Server is running.\nVisit /ws with a WebSocket client to test the live connection.")
	})

	// ---------- WebSocket route ----------
	// When someone visits /ws, Melody will upgrade the connection
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		// This line does the magic: upgrades HTTP → WebSocket
		err := m.HandleRequest(w, r)
		if err != nil {
			fmt.Println("WebSocket upgrade failed:", err)
		}
	})

	// ---------- WebSocket event handlers ----------
	// Runs once when a new client connects
	m.HandleConnect(func(s *melody.Session) {
		fmt.Printf("New client connected! Total clients: %d\n", m.Len())
		// Send a welcome message to ONLY this new client
		s.Write([]byte(`{"type":"welcome","message":"You are connected to the music queue!"}`))
	})

	// Runs when a client disconnects
	m.HandleDisconnect(func(s *melody.Session) {
		fmt.Printf("Client disconnected. Total clients now: %d\n", m.Len())
	})

	// Runs when we receive a message from any client
	m.HandleMessage(func(s *melody.Session, msg []byte) {
		fmt.Printf("Received message: %s\n", msg)
		// Echo it back to everyone (including sender) — just to prove it works
		m.Broadcast(msg)
	})

	// ---------- Start the server ----------
	fmt.Println("Server starting on http://localhost:8080")
	fmt.Println("• Open http://localhost:8080 in browser for hello page")
	fmt.Println("• Use a WebSocket tester for ws://localhost:8080/ws")
	
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Failed to start server:", err)
	}
}