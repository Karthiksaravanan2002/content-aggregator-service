package helper

import (
	"math"
	"net/http"

	"dev.azure.com/daimler-mic/content-aggregator/service/errors"

	"dev.azure.com/daimler-mic/content-aggregator/service/models"
)

func SelectRespStatusCode(resp *models.AggregateResponse) int {

	hasData := false
	hasError := false

	for _, providerResp := range resp.Providers {
		if len(providerResp.Data) > 0 {
			hasData = true
		}
		if len(providerResp.FeatureErrors) > 0 {
			hasError = true
		}
	}

	// Case 1: All providers succeeded (no errors anywhere)
	if !hasError {
		return http.StatusOK // 200
	}

	// Case 2: Mixed success + failure → 207 Multi-Status
	if hasData && hasError {
		return http.StatusMultiStatus // Multi-Status
	}

	// Case 3: Everything failed → choose highest priority error
	return SelectPriorityError(resp)

}

func SelectPriorityError(resp *models.AggregateResponse) int {

	var bestErr errors.AppError
	bestPriority := math.MaxInt32
	for _, prov := range resp.Providers {

		// Check for feature errors
		for _, ferr := range prov.FeatureErrors {
			if ferr == nil {
				continue
			}

			priority := errors.ErrorPriority(ferr.StatusCode())

			if priority < bestPriority {
				bestPriority = priority
				bestErr = ferr
			}
		}
	}

	if bestErr == nil {
		return http.StatusInternalServerError
	}

	return bestErr.StatusCode()
}
