// THIS FILE IS AUTOMATICALLY GENERATED. DO NOT EDIT.

// Package mobileanalyticsiface provides an interface for the Amazon Mobile Analytics.
package mobileanalyticsiface

import (
	"github.com/aws/aws-sdk-go/service/mobileanalytics"
)

// MobileAnalyticsAPI is the interface type for mobileanalytics.MobileAnalytics.
type MobileAnalyticsAPI interface {
	PutEvents(*mobileanalytics.PutEventsInput) (*mobileanalytics.PutEventsOutput, error)
}