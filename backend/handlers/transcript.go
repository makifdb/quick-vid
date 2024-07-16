package handlers

import (
	"net/http"

	"github.com/makifdb/quick-vid/services"
	"github.com/makifdb/quick-vid/utils"
)

// TranscriptHandler handles the /api/transcript/ endpoint
func TranscriptHandler(w http.ResponseWriter, r *http.Request) {

	// Extract videoId from the URL path
	videoId := r.PathValue("id")
	if videoId == "" {
		utils.RespondWithError(w, http.StatusBadRequest, "videoId not provided")
		return
	}

	utils.LogInfo("Request for videoId: %s", videoId)

	// Validate videoId
	if valid, err := services.ValidateID(videoId); !valid || err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Invalid videoId")
		return
	}

	// Get the transcript
	transcript, err := services.GetTranscript(videoId)
	if err != nil {
		utils.RespondWithError(w, http.StatusNotFound, "Transcript not found")
		return
	}

	utils.LogInfo("Transcript found for videoId: %s", videoId)
	utils.LogInfo("Transcript length: %d", len(*transcript))

	res, err := services.ProcessTranscript(transcript)
	if err != nil {
		utils.RespondWithError(w, http.StatusInternalServerError, "Failed to process transcript")
		return
	}

	// Return the transcript as JSON
	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"summary": *res})
}
