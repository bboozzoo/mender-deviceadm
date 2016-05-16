// Copyright 2016 Mender Software AS
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.

package main

import (
	"errors"
)

var (
	// device not found
	ErrDevNotFound = errors.New("not found")
)

type DataStore interface {
	GetDevices(skip, limit int, status string) ([]Device, error)
	// find a device with given `id`, returns the device or nil,
	// if device was not found, error is set to ErrDevNotFound
	GetDevice(id DeviceID) (*Device, error)
	// update device information
	PutDevice(dev *Device) error
}
