package packer_azure_image_version

import (
	"github.com/Azure/azure-sdk-for-go/profiles/latest/compute/mgmt/compute"
	"github.com/Azure/go-autorest/autorest"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"github.com/sirupsen/logrus"
)

type session struct {
	authorizer                 *autorest.Authorizer
	galleryImageVersionsClient map[string]*compute.GalleryImageVersionsClient
	galleryImagesClient        map[string]*compute.GalleryImagesClient
}

func (s *session) getAuthorizer() error {
	if s.authorizer != nil {
		return nil
	}
	// try from environment first
	a, err := auth.NewAuthorizerFromEnvironment()
	if err == nil {
		s.authorizer = &a

		logrus.Debug("retrieved authorizer from environment")

		return nil
	}

	a, err = auth.NewAuthorizerFromCLI()
	if err == nil {
		s.authorizer = &a

		logrus.Debug("retrieved authorizer from cli")

		return nil
	}

	return err
}

//
// // getGalleryImageClient creates a new GalleryImageClient instance and stores it in the provided session.
// // if an authorizer instance is missing, it will make a call to create it and then store in the session also.
// func (s *session) getGalleryImageClient(subID string) (err error) {
// 	if s.galleryImagesClient == nil {
// 		s.galleryImagesClient = make(map[string]*compute.GalleryImagesClient)
// 	}
//
// 	if s.galleryImagesClient[subID] != nil {
// 		logrus.Debugf("re-using gallery image client for subscription: %s", subID)
//
// 		return nil
// 	}
//
// 	logrus.Debugf("creating gallery image client for subscription: %s", subID)
//
// 	c := compute.NewGalleryImagesClient(subID)
// 	err = s.getAuthorizer()
// 	if err != nil {
// 		return
// 	}
//
// 	s.galleryImagesClient[subID] = &c
// 	s.galleryImagesClient[subID].Authorizer = *s.authorizer
//
// 	return
// }

// getGalleryImageVersionsClient creates a new GalleryImageVersionsClient instance and stores it in the provided session.
// if an authorizer instance is missing, it will make a call to create it and then store in the session also.
func (s *session) getGalleryImageVersionsClient(subID string) (err error) {
	if s.galleryImageVersionsClient == nil {
		s.galleryImageVersionsClient = make(map[string]*compute.GalleryImageVersionsClient)
	}

	if s.galleryImageVersionsClient[subID] != nil {
		logrus.Debugf("re-using gallery image versions client for subscription: %s", subID)

		return nil
	}

	logrus.Debugf("creating gallery image versions client for subscription: %s", subID)

	c := compute.NewGalleryImageVersionsClient(subID)
	err = s.getAuthorizer()
	if err != nil {
		return
	}

	s.galleryImageVersionsClient[subID] = &c
	s.galleryImageVersionsClient[subID].Authorizer = *s.authorizer

	return
}
