/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

import (
	arbv1 "github.com/kubernetes-incubator/kube-arbitrator/pkg/apis/v1"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// QueueLister helps list Queues.
type QueueLister interface {
	// List lists all Queues in the indexer.
	List(selector labels.Selector) (ret []*arbv1.Queue, err error)
	// Queues returns an object that can list and get Queues.
	Queues(namespace string) QueueNamespaceLister
}

// queueLister implements the QueueLister interface.
type queueLister struct {
	indexer cache.Indexer
}

// NewQueueLister returns a new QueueLister.
func NewQueueLister(indexer cache.Indexer) QueueLister {
	return &queueLister{indexer: indexer}
}

// List lists all Queues in the indexer.
func (s *queueLister) List(selector labels.Selector) (ret []*arbv1.Queue, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*arbv1.Queue))
	})
	return ret, err
}

// Queues returns an object that can list and get Queues.
func (s *queueLister) Queues(namespace string) QueueNamespaceLister {
	return queueNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// QueueNamespaceLister helps list and get Queues.
type QueueNamespaceLister interface {
	// List lists all Queues in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*arbv1.Queue, err error)
	// Get retrieves the Queue from the indexer for a given namespace and name.
	Get(name string) (*arbv1.Queue, error)
}

// queueNamespaceLister implements the QueueNamespaceLister
// interface.
type queueNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all Queues in the indexer for a given namespace.
func (s queueNamespaceLister) List(selector labels.Selector) (ret []*arbv1.Queue, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*arbv1.Queue))
	})
	return ret, err
}

// Get retrieves the Queue from the indexer for a given namespace and name.
func (s queueNamespaceLister) Get(name string) (*arbv1.Queue, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(arbv1.Resource("queue"), name)
	}
	return obj.(*arbv1.Queue), nil
}
