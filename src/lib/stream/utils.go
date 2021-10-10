package stream

import "errors"

func validateOpts(opts Config) error {
	if AWSRegion == "" {
		return errors.New("AWSRegion is required")
	}

	if URL == "" {
		return errors.New("A valid SQS URL is required")
	}

	if opts.BatchSize < 0 || opts.BatchSize > 10 {
		return errors.New("BatchSize should be between 1-10")
	}

	if opts.WaitSeconds < 0 || opts.WaitSeconds > 20 {
		return errors.New("WaitSecond should be between 1-20")
	}

	if opts.VisibilityTimeout < 0 || opts.VisibilityTimeout > 12*60*60 {
		return errors.New("WaitSecond should be between 1-43200")
	}

	return nil
}

