package try

type Try func() error

func RunOrError(calls ...Try) error {
	for _, c := range calls {
		if err := c(); err != nil {
			return err
		}
	}
	return nil
}