package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/lulzshadowwalker/pupsik/database"
	"github.com/lulzshadowwalker/pupsik/types"
	"github.com/lulzshadowwalker/pupsik/utils"
	"github.com/lulzshadowwalker/pupsik/view/home"
	"github.com/replicate/replicate-go"
)

type (
	ReplicateResponse struct {
		Input struct {
			Prompt string `json:"prompt"`
		} `json:"input"`
		Status ReplicateResponseStatus `json:"status"`
		Output []string                `json:"output"`
	}

	ReplicateResponseStatus string
)

const (
	ReplicateResponseStatusSucceeded ReplicateResponseStatus = "succeeded"
	ReplicateResponseStatusFailed    ReplicateResponseStatus = "failed"
)

func HandleGalleryCreate(w http.ResponseWriter, r *http.Request) error {
	user, err := utils.GetUserFromContext(r.Context())
	if err != nil {
		return err
	}

	log.Println("prompt", r.FormValue("prompt"))

	// TODO: handle errors e.g. form value not found
	img := types.Image{
		UserID:  user.ID,
		BatchID: uuid.New(),
		Prompt:  r.FormValue("prompt"),
		Status:  types.ImageStatusPending,
	}

	img, err = database.StoreImage(r.Context(), img, nil)
	if err != nil {
		return err
	}

	if err = generateImage(r.Context(), img.Prompt, 1, user.ID, img.BatchID); err != nil {
		return err
	}

	return render(w, r, home.Image(img))
}

// TODO: this most likely needs rework im sleepy af
func HandleGalleryImageStatus(w http.ResponseWriter, r *http.Request) error {
	user, err := utils.GetUserFromContext(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return err
	}

	idParam := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return errors.New("bad request")
	}

	img, err := database.GetImageByID(r.Context(), id, nil)
	if err != nil {
		return err
	}

	if img.UserID != user.ID {
		w.WriteHeader(http.StatusForbidden)
		return errors.New("forbidden")
	}

	return nil
}

// TODO: handle count being more than one properly
func generateImage(ctx context.Context, prompt string, count int, userID, batchID uuid.UUID) error {
	r8, err := replicate.NewClient(replicate.WithTokenFromEnv())
	if err != nil {
		return fmt.Errorf("failed to init replicate client because %w", err)
	}

	latestConsistencyModel := "553803fd018b3cf875a8bc774c99da9b33f36647badfd88a6eec90d61c5f62fc"
	input := replicate.PredictionInput{
		"prompt":      prompt,
		"num_outputs": count,
	}

	// callbackURL, err := utils.GetURL(fmt.Sprintf("/generate/callback/%d/%d", userID, batchID))
	// if err != nil {
	// 	return err
	// }

	callbackURL := fmt.Sprintf("https://webhook.site/31a872ed-cde3-4946-bea5-23358a62640c/%s/%s", userID, batchID)
	webhook := replicate.Webhook{
		URL:    callbackURL,
		Events: []replicate.WebhookEventType{"completed"},
	}

	_, err = r8.CreatePrediction(ctx, latestConsistencyModel, input, &webhook, false)
	if err != nil {
		return err
	}

	return nil
}

func HandleReplicateCallback(_ http.ResponseWriter, r *http.Request) error {
	log.Println("REPLICATE CALLBACK")
	var response ReplicateResponse
	if err := json.NewDecoder(r.Body).Decode(&response); err != nil {
		return err
	}

	switch response.Status {
	case ReplicateResponseStatusSucceeded:
	case ReplicateResponseStatusFailed:
		// TODO: handle failed replicate
	default:
		return fmt.Errorf("unrecognized status %q", response.Status)
	}

	batchID, err := uuid.Parse(chi.URLParam(r, "batchID"))
	if err != nil {
		return fmt.Errorf("replicate callback invalid batchID %q", err)
	}

	images, err := database.GetImagesByBatchID(r.Context(), batchID, nil)
	if err != nil {
		return fmt.Errorf("replicate callback failed to find image with batchID %s because %w", batchID, err)
	}

	if len(images) != len(response.Output) {
		return fmt.Errorf("replicate callback un-equal images compaired to replicate outputs got %d expected %d", len(response.Output), len(images))
	}

	tx, err := database.DB.BeginTx(r.Context(), nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction because %w", err)
	}

	for i, imageURL := range response.Output {
		images[i].Status = types.ImageStatusFinished

		// TODO: Store image to supabase
		images[i].URL = imageURL
		images[i].Prompt = response.Input.Prompt
		if _, err := database.UpdateImage(r.Context(), images[i], tx); err != nil {
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction because %w", err)
	}

	return nil
}

func HandleGenerateImageStatus(w http.ResponseWriter, r *http.Request) error {
	// TODO: handle auth and owner user

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		return err
	}
	image, err := database.GetImageByID(r.Context(), id, nil)
	if err != nil {
		return err
	}

	slog.Info("checking image status", "id", id)
	return render(w, r, home.Image(image))
}
