package fakes

import (
	"github.com/frodenas/bosh-openstack-cpi/openstack/snapshot_service"
)

type FakeSnapshotService struct {
	CreateCalled      bool
	CreateErr         error
	CreateID          string
	CreateVolumeID    string
	CreateDescription string

	DeleteCalled bool
	DeleteErr    error

	FindCalled   bool
	FindFound    bool
	FindSnapshot snapshot.Snapshot
	FindErr      error
}

func (s *FakeSnapshotService) Create(volumeID string, description string) (string, error) {
	s.CreateCalled = true
	s.CreateVolumeID = volumeID
	s.CreateDescription = description
	return s.CreateID, s.CreateErr
}

func (s *FakeSnapshotService) Delete(id string) error {
	s.DeleteCalled = true
	return s.DeleteErr
}

func (s *FakeSnapshotService) Find(id string) (snapshot.Snapshot, bool, error) {
	s.FindCalled = true
	return s.FindSnapshot, s.FindFound, s.FindErr
}
