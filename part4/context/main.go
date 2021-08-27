package main

import (
	"context"
	"fmt"
	"time"
)

func main()  {
	//var wg sync.WaitGroup
	//ctx, cancel := context.WithCancel(context.Background())
	//defer cancel()
	//
	//wg.Add(1)
	//go func() {
	//	defer wg.Done()
	//	if err := printGreeting(ctx); err != nil {
	//		fmt.Printf("cannot print greeting: %v \n", err)
	//		cancel()
	//	}
	//}()
	//
	//wg.Add(1)
	//go func() {
	//	defer wg.Done()
	//	if err := printFarewell(ctx); err != nil {
	//		fmt.Printf("cannot print farewell: %v \n", err)
	//	}
	//}()
	//
	//wg.Wait()
	ProcessRequest("jane", "abc123")
}

func printGreeting(ctx context.Context) error  {
	greeting, err := genGreeting(ctx)
	if err != nil {
		return err
	}
	fmt.Printf("%s world! \n", greeting)
	return nil
}

func printFarewell(ctx context.Context) error {
	farewell, err := genFarewell(ctx)
	if err != nil {
		return err
	}
	fmt.Printf("%s world! \n", farewell)
	return nil
}

func genGreeting(ctx context.Context) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	switch locale, err  := locale(ctx) ;{
	case err != nil:
		return "", ctx.Err()
	case locale == "EN/US":
		return "hello", nil
	}
	return "", fmt.Errorf("unsupoprt locale")
}

func genFarewell(ctx context.Context) (string, error)  {
	switch locale, err := locale(ctx); {
	case err != nil:
		return "", err
	case locale == "EN/US":
		return "goodbye", nil
	}
	return "", fmt.Errorf("unsupported locale")
}

func locale(ctx context.Context) (string, error)  {
	if deadline, ok := ctx.Deadline(); ok {
		if deadline.Sub(time.Now().Add(1*time.Minute)) <= 0 {
			return "", context.DeadlineExceeded
		}
	}

	select {
	case <-ctx.Done():
		return "", fmt.Errorf("canceled")
	case <-time.After(1*time.Minute):
	}
	return "EN/US", nil
}

type ctxKey int

const (
	ctxUserID ctxKey = iota
	ctxAuthToken
)

func UserID(c context.Context) string {
	return c.Value(ctxUserID).(string)
}

func AuthToken(c context.Context) string {
	return c.Value(ctxAuthToken).(string)
}

func ProcessRequest(userID, authToken string)  {
	ctx := context.WithValue(context.Background(), ctxUserID, userID)
	ctx = context.WithValue(ctx, ctxAuthToken, authToken)
	HandleResponse(ctx)
}

func HandleResponse(ctx context.Context) {
	fmt.Printf(
		"handling response for %v (auth: %v)",
		UserID(ctx),
		AuthToken(ctx),
	)
}
