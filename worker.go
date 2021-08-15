func (c *Controller) runWorker(ctx context.Context) {
	for c.processNextItem(ctx) {
	}
}

func (c *Controller) processNextItem(ctx context.Context) bool {
	obj, shutdown := c.queue.Get()
	if shutdown {
		return false
	}
	defer c.queue.Done(obj)

	err := c.processEvent(ctx, obj)
	if err == nil {
		c.logger.Debug("processed item")
		c.queue.Forget(obj)
	} else if c.queue.NumRequeues(obj) < maxRetries {
		c.logger.Errorf("error processing event: %v, retrying", err)
		c.queue.AddRateLimited(obj)
	} else {
		c.logger.Errorf("error processing event: %v, max retries reached", err)
		c.queue.Forget(obj)
		utilruntime.HandleError(err)
	}

	return true
}

func (c *Controller) processEvent(ctx context.Context, obj interface{}) error {
	event, ok := obj.(event)
	if !ok {
		c.logger.Error("unexpected event ", obj)
		return nil
	}
	switch event.eventType {
	case addEcho:
		return c.processAddEcho(ctx, event.newObj.(*echov1alpha1.Echo))
	case addScheduledEcho:
		return c.processAddScheduledEcho(ctx, event.newObj.(*echov1alpha1.ScheduledEcho))
	case updateScheduledEcho:
		return c.processUpdateScheduledEcho(
			ctx,
			event.oldObj.(*echov1alpha1.ScheduledEcho),
			event.newObj.(*echov1alpha1.ScheduledEcho),
		)
	}
	return nil
}

func (c *Controller) processAddEcho(ctx context.Context, echo *echov1alpha1.Echo) error {
	job := createJob(echo, c.namespace)
	exists, err := resourceExists(job, c.jobInformer.GetIndexer())
	if err != nil {
		return fmt.Errorf("error checking job existence %v", err)
	}
	if exists {
		c.logger.Debug("job already exists, skipping")
		return nil
	}

	_, err = c.kubeClientSet.BatchV1().
		Jobs(c.namespace).
		Create(ctx, job, metav1.CreateOptions{})
	return err
}

func resourceExists(obj interface{}, indexer cache.Indexer) (bool, error) {
	key, err := cache.MetaNamespaceKeyFunc(obj)
	if err != nil {
		return false, fmt.Errorf("error getting key %v", err)
	}
	_, exists, err := indexer.GetByKey(key)
	return exists, err
}