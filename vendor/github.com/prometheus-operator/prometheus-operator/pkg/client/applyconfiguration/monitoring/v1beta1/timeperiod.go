// Copyright The prometheus-operator Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1beta1

import (
	monitoringv1beta1 "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1beta1"
)

// TimePeriodApplyConfiguration represents a declarative configuration of the TimePeriod type for use
// with apply.
type TimePeriodApplyConfiguration struct {
	Times       []TimeRangeApplyConfiguration       `json:"times,omitempty"`
	Weekdays    []monitoringv1beta1.WeekdayRange    `json:"weekdays,omitempty"`
	DaysOfMonth []DayOfMonthRangeApplyConfiguration `json:"daysOfMonth,omitempty"`
	Months      []monitoringv1beta1.MonthRange      `json:"months,omitempty"`
	Years       []monitoringv1beta1.YearRange       `json:"years,omitempty"`
}

// TimePeriodApplyConfiguration constructs a declarative configuration of the TimePeriod type for use with
// apply.
func TimePeriod() *TimePeriodApplyConfiguration {
	return &TimePeriodApplyConfiguration{}
}

// WithTimes adds the given value to the Times field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the Times field.
func (b *TimePeriodApplyConfiguration) WithTimes(values ...*TimeRangeApplyConfiguration) *TimePeriodApplyConfiguration {
	for i := range values {
		if values[i] == nil {
			panic("nil value passed to WithTimes")
		}
		b.Times = append(b.Times, *values[i])
	}
	return b
}

// WithWeekdays adds the given value to the Weekdays field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the Weekdays field.
func (b *TimePeriodApplyConfiguration) WithWeekdays(values ...monitoringv1beta1.WeekdayRange) *TimePeriodApplyConfiguration {
	for i := range values {
		b.Weekdays = append(b.Weekdays, values[i])
	}
	return b
}

// WithDaysOfMonth adds the given value to the DaysOfMonth field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the DaysOfMonth field.
func (b *TimePeriodApplyConfiguration) WithDaysOfMonth(values ...*DayOfMonthRangeApplyConfiguration) *TimePeriodApplyConfiguration {
	for i := range values {
		if values[i] == nil {
			panic("nil value passed to WithDaysOfMonth")
		}
		b.DaysOfMonth = append(b.DaysOfMonth, *values[i])
	}
	return b
}

// WithMonths adds the given value to the Months field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the Months field.
func (b *TimePeriodApplyConfiguration) WithMonths(values ...monitoringv1beta1.MonthRange) *TimePeriodApplyConfiguration {
	for i := range values {
		b.Months = append(b.Months, values[i])
	}
	return b
}

// WithYears adds the given value to the Years field in the declarative configuration
// and returns the receiver, so that objects can be build by chaining "With" function invocations.
// If called multiple times, values provided by each call will be appended to the Years field.
func (b *TimePeriodApplyConfiguration) WithYears(values ...monitoringv1beta1.YearRange) *TimePeriodApplyConfiguration {
	for i := range values {
		b.Years = append(b.Years, values[i])
	}
	return b
}
