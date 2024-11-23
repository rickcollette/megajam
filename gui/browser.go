package gui

import (
    "bytes"
    "fmt"
    "image"
    "log"
    "os"
    "strings"
    "sync"

    "megajam/playlist"

    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/canvas"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/dialog"
    "fyne.io/fyne/v2/widget"
    "github.com/dhowden/tag"
    "github.com/gopxl/beep/v2/mp3"
)

// ExtractAlbumArt extracts album art from an MP3 file.
func ExtractAlbumArt(filePath string) *canvas.Image {
    f, err := os.Open(filePath)
    if err != nil {
        log.Printf("Error opening file %s: %v", filePath, err)
        return nil
    }
    defer f.Close()

    meta, err := tag.ReadFrom(f)
    if err != nil {
        log.Printf("Error reading tags from file %s: %v", filePath, err)
        return nil
    }

    picture := meta.Picture()
    if picture == nil {
        log.Printf("No album art found in file %s", filePath)
        return nil
    }

    img, _, err := image.Decode(bytes.NewReader(picture.Data))
    if err != nil {
        log.Printf("Error decoding album art from file %s: %v", filePath, err)
        return nil
    }

    canvasImage := canvas.NewImageFromImage(img)
    canvasImage.SetMinSize(fyne.NewSize(100, 100)) // Set a fixed size for thumbnails
    return canvasImage
}

// ExtractMetadata extracts metadata (Title, Artist, Length) from an MP3 file.
func ExtractMetadata(filePath string) (title, artist, length string) {
    f, err := os.Open(filePath)
    if err != nil {
        log.Printf("Error opening file %s: %v", filePath, err)
        return "Unknown", "Unknown", "Unknown"
    }
    defer f.Close()

    meta, err := tag.ReadFrom(f)
    if err != nil {
        log.Printf("Error reading tags from file %s: %v", filePath, err)
        return "Unknown", "Unknown", "Unknown"
    }

    title = meta.Title()
    if title == "" {
        title = "Unknown"
    }

    artist = meta.Artist()
    if artist == "" {
        artist = "Unknown"
    }

    // Attempt to get track length using Beep's mp3 decoder
    streamer, format, err := mp3.Decode(f)
    if err != nil {
        log.Printf("Error decoding MP3 for length extraction %s: %v", filePath, err)
        length = "Unknown"
    } else {
        defer streamer.Close()
        totalSamples := streamer.Len()
        if totalSamples == -1 {
            length = "Unknown"
        } else {
            duration := float64(totalSamples) / float64(format.SampleRate)
            minutes := int(duration) / 60
            seconds := int(duration) % 60
            length = fmt.Sprintf("%d:%02d", minutes, seconds)
        }
    }

    return title, artist, length
}

// createEnhancedBrowserSection creates the playlist browser with enhanced features.
// It returns the container and a channel to send the selected track index.
func createEnhancedBrowserSection(playlist *playlist.Playlist, addTrackButton, removeTrackButton *widget.Button, myApp fyne.App, myWindow fyne.Window) (*fyne.Container, chan int) {
    // Sidebar navigation
    sidebar := container.NewVBox(
        widget.NewButton("Deezer", func() {
            dialog.ShowInformation("Info", "Deezer integration coming soon!", myWindow)
        }),
        widget.NewButton("TIDAL", func() {
            dialog.ShowInformation("Info", "TIDAL integration coming soon!", myWindow)
        }),
        widget.NewButton("Beatport", func() {
            dialog.ShowInformation("Info", "Beatport integration coming soon!", myWindow)
        }),
        widget.NewButton("SoundCloud", func() {
            dialog.ShowInformation("Info", "SoundCloud integration coming soon!", myWindow)
        }),
        widget.NewButton("Offline Cache", func() {
            dialog.ShowInformation("Info", "Offline Cache feature coming soon!", myWindow)
        }),
    )

    // Search bar
    searchEntry := widget.NewEntry()
    searchEntry.SetPlaceHolder("Search...")
    var filteredTracks []string
    var mutex sync.Mutex

    selectedTrackChan := make(chan int)

    updateFilteredTracks := func(query string) {
        mutex.Lock()
        defer mutex.Unlock()
        filteredTracks = []string{}
        for _, track := range playlist.Tracks {
            title, artist, _ := ExtractMetadata(track)
            if strings.Contains(strings.ToLower(title), strings.ToLower(query)) ||
                strings.Contains(strings.ToLower(artist), strings.ToLower(query)) {
                filteredTracks = append(filteredTracks, track)
            }
        }
        trackList.Refresh()
    }

    searchEntry.OnChanged = updateFilteredTracks

    // Thumbnail generation (asynchronous)
    thumbnailContainer := container.NewHBox()
    var wg sync.WaitGroup
    for _, trackPath := range playlist.Tracks {
        wg.Add(1)
        go func(tp string) {
            defer wg.Done()
            thumbnail := ExtractAlbumArt(tp)
            if thumbnail != nil {
                myApp.QueueUpdate(func() {
                    thumbnailContainer.Add(thumbnail)
                    thumbnailContainer.Refresh()
                })
            }
        }(trackPath)
    }
    wg.Wait()
    thumbnailScroll := container.NewHScroll(thumbnailContainer)

    // Column headers
    columnHeaders := container.NewGridWithColumns(3,
        widget.NewLabelWithStyle("Title", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
        widget.NewLabelWithStyle("Artist", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
        widget.NewLabelWithStyle("Length", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
    )

    // Track list with metadata using widget.List for selection
    var selectedTrackIndex int = -1
    trackList := widget.NewList(
        func() int {
            mutex.Lock()
            defer mutex.Unlock()
            if len(filteredTracks) > 0 {
                return len(filteredTracks)
            }
            return len(playlist.Tracks)
        },
        func() fyne.CanvasObject {
            return container.NewHBox(
                widget.NewLabel("Title"),
                widget.NewLabel("Artist"),
                widget.NewLabel("Length"),
            )
        },
        func(i widget.ListItemID, o fyne.CanvasObject) {
            mutex.Lock()
            defer mutex.Unlock()
            var trackPath string
            if len(filteredTracks) > 0 {
                trackPath = filteredTracks[i]
            } else {
                trackPath = playlist.Tracks[i]
            }
            title, artist, length := ExtractMetadata(trackPath)
            labels := o.(*fyne.Container).Objects
            labels[0].(*widget.Label).SetText(title)
            labels[1].(*widget.Label).SetText(artist)
            labels[2].(*widget.Label).SetText(length)
        },
    )

    trackList.OnSelected = func(id widget.ListItemID) {
        mutex.Lock()
        defer mutex.Unlock()
        if len(filteredTracks) > 0 {
            selectedTrackIndex = findIndex(playlist.Tracks, filteredTracks[id])
        } else {
            selectedTrackIndex = id
        }
        log.Printf("Selected track index: %d", selectedTrackIndex)
        selectedTrackChan <- selectedTrackIndex
    }

    // Helper function to find the index of a track in the main playlist
    findIndex := func(tracks []string, target string) int {
        for i, track := range tracks {
            if track == target {
                return i
            }
        }
        return -1
    }

    // Assemble browser layout
    return container.NewBorder(
        container.NewVBox(searchEntry, thumbnailScroll, columnHeaders), // Top: Search, thumbnails, and column headers
        container.NewHBox(addTrackButton, removeTrackButton),           // Bottom: Add/Remove buttons
        sidebar,                                                        // Left: Sidebar
        nil,                                                            // Right: None
        trackList,                                                      // Center: Track list
    ), selectedTrackChan
}