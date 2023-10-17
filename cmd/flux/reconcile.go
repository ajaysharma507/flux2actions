/*
Copyright 2020 The Flux authors

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

package main

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
	apimeta "k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/util/retry"
	kstatus "sigs.k8s.io/cli-utils/pkg/kstatus/status"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/fluxcd/pkg/apis/meta"

	"github.com/fluxcd/flux2/v2/internal/utils"
)

var reconcileCmd = &cobra.Command{
	Use:   "reconcile",
	Short: "Reconcile sources and resources",
	Long:  `The reconcile sub-commands trigger a reconciliation of sources and resources.`,
}

func init() {
	rootCmd.AddCommand(reconcileCmd)
}

type reconcileCommand struct {
	apiType
	object reconcilable
}

type reconcilable interface {
	adapter     // to be able to load from the cluster
	copyable    // to be able to calculate patches
	suspendable // to tell if it's suspended

	// these are implemented by anything embedding metav1.ObjectMeta
	GetAnnotations() map[string]string
	SetAnnotations(map[string]string)

	hasReconciler() bool                 // does the object's current configuration has a reconciler to reconcile it?
	lastHandledReconcileRequest() string // what was the last handled reconcile request?
	successMessage() string              // what do you want to tell people when successfully reconciled?
}

func reconcilableConditions(object reconcilable) []metav1.Condition {
	if s, ok := object.(meta.ObjectWithConditions); ok {
		return s.GetConditions()
	}

	if s, ok := object.(oldConditions); ok {
		return *s.GetStatusConditions()
	}

	return []metav1.Condition{}
}

func (reconcile reconcileCommand) run(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("%s name is required", reconcile.kind)
	}
	name := args[0]

	ctx, cancel := context.WithTimeout(context.Background(), rootArgs.timeout)
	defer cancel()

	kubeClient, err := utils.KubeClient(kubeconfigArgs, kubeclientOptions)
	if err != nil {
		return err
	}

	namespacedName := types.NamespacedName{
		Namespace: *kubeconfigArgs.Namespace,
		Name:      name,
	}

	err = kubeClient.Get(ctx, namespacedName, reconcile.object.asClientObject())
	if err != nil {
		return err
	}

	if !reconcile.object.hasReconciler() {
		logger.Successf("reconciliation not supported by the object")
		return nil
	}

	if reconcile.object.isSuspended() {
		return fmt.Errorf("resource is suspended")
	}

	logger.Actionf("annotating %s %s in %s namespace", reconcile.kind, name, *kubeconfigArgs.Namespace)
	if err := requestReconciliation(ctx, kubeClient, namespacedName,
		reconcile.groupVersion.WithKind(reconcile.kind)); err != nil {
		return err
	}
	logger.Successf("%s annotated", reconcile.kind)

	lastHandledReconcileAt := reconcile.object.lastHandledReconcileRequest()
	logger.Waitingf("waiting for %s reconciliation", reconcile.kind)
	if err := wait.PollUntilContextTimeout(ctx, rootArgs.pollInterval, rootArgs.timeout, true,
		reconciliationHandled(kubeClient, namespacedName, reconcile.object, lastHandledReconcileAt)); err != nil {
		return err
	}
	readyCond := apimeta.FindStatusCondition(reconcilableConditions(reconcile.object), meta.ReadyCondition)
	if readyCond == nil {
		return fmt.Errorf("status can't be determined")
	}

	if readyCond.Status != metav1.ConditionTrue {
		return fmt.Errorf("%s reconciliation failed: '%s'", reconcile.kind, readyCond.Message)
	}
	logger.Successf(reconcile.object.successMessage())
	return nil
}

func reconciliationHandled(kubeClient client.Client, namespacedName types.NamespacedName, obj reconcilable, lastHandledReconcileAt string) wait.ConditionWithContextFunc {
	return func(ctx context.Context) (bool, error) {
		err := kubeClient.Get(ctx, namespacedName, obj.asClientObject())
		if err != nil {
			return false, err
		}

		if obj.lastHandledReconcileRequest() == lastHandledReconcileAt {
			return false, nil
		}

		result, err := kstatusCompute(obj.asClientObject())
		if err != nil {
			return false, err
		}

		return result.Status == kstatus.CurrentStatus, nil
	}
}

func requestReconciliation(ctx context.Context, kubeClient client.Client,
	namespacedName types.NamespacedName, gvk schema.GroupVersionKind) error {
	return retry.RetryOnConflict(retry.DefaultBackoff, func() (err error) {
		object := &metav1.PartialObjectMetadata{}
		object.SetGroupVersionKind(gvk)
		object.SetName(namespacedName.Name)
		object.SetNamespace(namespacedName.Namespace)
		if err := kubeClient.Get(ctx, namespacedName, object); err != nil {
			return err
		}
		patch := client.MergeFrom(object.DeepCopy())
		if ann := object.GetAnnotations(); ann == nil {
			object.SetAnnotations(map[string]string{
				meta.ReconcileRequestAnnotation: time.Now().Format(time.RFC3339Nano),
			})
		} else {
			ann[meta.ReconcileRequestAnnotation] = time.Now().Format(time.RFC3339Nano)
			object.SetAnnotations(ann)
		}
		return kubeClient.Patch(ctx, object, patch)
	})
}
