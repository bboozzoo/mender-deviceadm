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
	"github.com/pkg/errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	DbName        = "deviceadm"
	DbDevicesColl = "devices"
)

type DataStoreMongo struct {
	session *mgo.Session
}

func NewDataStoreMongo(host string) (*DataStoreMongo, error) {
	s, err := mgo.Dial(host)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open mgo session")
	}
	return &DataStoreMongo{session: s}, nil
}

func (db *DataStoreMongo) GetDevices(skip, limit int, status string) ([]Device, error) {
	s := db.session.Copy()
	defer s.Close()
	c := s.DB(DbName).C(DbDevicesColl)
	res := []Device{}

	var filter bson.M
	if status != "" {
		filter = bson.M{"status": status}
	}

	err := c.Find(filter).Skip(skip).Limit(limit).All(&res)

	if err != nil {
		return nil, errors.Wrap(err, "failed to fetch device list")
	}

	return res, nil
}

func (db *DataStoreMongo) GetDevice(id DeviceID) (*Device, error) {
	s := db.session.Copy()
	defer s.Close()
	c := s.DB(DbName).C(DbDevicesColl)

	filter := bson.M{"id": id}
	res := Device{}

	err := c.Find(filter).One(&res)

	if err != nil {
		if err == mgo.ErrNotFound {
			return nil, ErrDevNotFound
		} else {
			return nil, errors.Wrap(err, "failed to fetch device")
		}
	}

	return &res, nil
}

func (db *DataStoreMongo) PutDevice(dev *Device) error {
	s := db.session.Copy()
	defer s.Close()
	c := s.DB(DbName).C(DbDevicesColl)

	filter := bson.M{"id": dev.ID}
	updates := map[string]interface{}{}
	if dev.Status != "" {
		updates["status"] = dev.Status
	}
	if dev.Key != "" {
		updates["key"] = dev.Key
	}
	if dev.DeviceIdentity != "" {
		updates["device_identity"] = dev.DeviceIdentity
	}
	// TODO: should attributes be merged?
	if len(dev.Attributes) != 0 {
		updates["attributes"] = dev.Attributes
	}

	data := bson.M{"$set": updates}
	_, err := c.Upsert(filter, data)
	if err != nil {
		return errors.Wrap(err, "failed to store device")
	}
	return nil
}
