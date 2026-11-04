# Pit Wall Fuel Planner

Fuel loads, stop counts and pit windows for Australian circuits, adjusted for the race-day forecast at the track.

Open `index.html` in any browser — it's a single self-contained page.

- 17 circuits, including One Raceway (Goulburn NSW)
- Car classes from Supercars and GT3 down to Toyota 86 / GR86, Subaru BRZ and club cars
- Fuel use, lap times and pit losses are class estimates; replace them with your own logged figures
- The forecast is a snapshot (Open-Meteo, 4 Nov 2026). "Refresh forecast" pulls the latest Open-Meteo forecast for the selected circuit (needs an internet connection).

## Windows app

`app/` holds a small Go launcher that bundles the planner into a single `FuelPlanner.exe`.
Double-clicking it opens the planner in its own Edge (or Chrome) window. Nothing to install.
The app closes itself about 20 s after the window is closed.

Build (from Linux, macOS or Windows with Go installed):

    cp index.html app/index.html
    cd app
    GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -ldflags "-s -w -H=windowsgui" -o FuelPlanner.exe .
