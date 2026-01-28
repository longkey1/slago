package cmd

import (
	"fmt"
	"sync"
	"time"

	"github.com/longkey1/slago/internal/collector"
	"github.com/longkey1/slago/internal/config"
	"github.com/longkey1/slago/internal/output"
	"github.com/longkey1/slago/internal/slack"
	"github.com/longkey1/slago/internal/dateutil"
	"github.com/spf13/cobra"
)

var (
	listDay             string
	listMonth           string
	listFrom            string
	listTo              string
	listThread          bool
	listAuthor          string
	listMentions        []string
	listChannels        []string
	listExcludeChannels []string
	listParallel        int
)

func newListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "Collect messages for a date range and save to files",
		Long: `Collect Slack messages for a date range and save to JSON files.

Output is saved to logs/YYYY/MM/DD/slack.json for each day.

Date range options (mutually exclusive):
  --day      Single day (YYYY-MM-DD)
  --month    Entire month (YYYY-MM)
  --from/--to Custom range (both required)

Examples:
  slago list --day 2025-01-15
  slago list --month 2025-01
  slago list --from 2025-01-01 --to 2025-01-15
  slago list -m 2025-01 --thread --author U12345678
  slago list -d 2025-01-15 --mention U111 --mention @team
  slago list -m 2025-01 --channel general --channel random
  slago list -d 2025-01-15 --exclude-channel announcements`,
		RunE: runList,
	}

	cmd.Flags().StringVarP(&listDay, "day", "d", "", "Day to collect (YYYY-MM-DD)")
	cmd.Flags().StringVarP(&listMonth, "month", "m", "", "Month to collect (YYYY-MM)")
	cmd.Flags().StringVar(&listFrom, "from", "", "Start date (YYYY-MM-DD)")
	cmd.Flags().StringVar(&listTo, "to", "", "End date (YYYY-MM-DD)")
	cmd.Flags().BoolVar(&listThread, "thread", false, "Get entire threads")
	cmd.Flags().StringVar(&listAuthor, "author", "", "Filter by author")
	cmd.Flags().StringSliceVar(&listMentions, "mention", nil, "Filter by mention (comma-separated User IDs or @group-names)")
	cmd.Flags().StringSliceVar(&listChannels, "channel", nil, "Filter by channel (comma-separated channel names)")
	cmd.Flags().StringSliceVar(&listExcludeChannels, "exclude-channel", nil, "Exclude channels (comma-separated channel names)")
	cmd.Flags().IntVarP(&listParallel, "parallel", "p", 1, "Number of parallel workers")

	return cmd
}

func runList(cmd *cobra.Command, args []string) error {
	// Load config
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Override from flags
	if token != "" {
		cfg.Token = token
	}
	if listAuthor == "" {
		listAuthor = cfg.Author
	}
	if len(listMentions) == 0 {
		listMentions = cfg.Mention
	}

	if err := cfg.Validate(); err != nil {
		return err
	}

	// Parse date range
	dateRange, err := parseDateRange()
	if err != nil {
		return err
	}

	// Create Slack client
	client := slack.NewClient(cfg.Token)

	// Get all days to process
	days := dateRange.Days()
	if len(days) == 0 {
		return fmt.Errorf("no days to process")
	}

	fmt.Printf("Collecting messages for %d day(s)...\n", len(days))

	// Process days with parallelism
	results := processdays(client, days, listParallel)

	// Report results
	var errors []error
	for _, result := range results {
		if result.Error != nil {
			errors = append(errors, fmt.Errorf("%s: %w", dateutil.FormatDate(result.Date), result.Error))
		} else {
			fmt.Printf("[INFO] %s: %d threads collected, saved to %s\n",
				dateutil.FormatDate(result.Date),
				len(result.Threads),
				dateutil.OutputPath(result.Date))
		}
	}

	if len(errors) > 0 {
		for _, err := range errors {
			fmt.Printf("[ERROR] %v\n", err)
		}
		return fmt.Errorf("%d day(s) failed", len(errors))
	}

	return nil
}

func parseDateRange() (dateutil.DateRange, error) {
	// Count how many date options are specified
	count := 0
	if listDay != "" {
		count++
	}
	if listMonth != "" {
		count++
	}
	if listFrom != "" || listTo != "" {
		count++
	}

	if count == 0 {
		return dateutil.DateRange{}, fmt.Errorf("date range required: use --day, --month, or --from/--to")
	}
	if count > 1 {
		return dateutil.DateRange{}, fmt.Errorf("only one date range option allowed: --day, --month, or --from/--to")
	}

	if listDay != "" {
		day, err := dateutil.ParseDay(listDay)
		if err != nil {
			return dateutil.DateRange{}, err
		}
		return dateutil.DayRange(day), nil
	}

	if listMonth != "" {
		return dateutil.ParseMonth(listMonth)
	}

	if listFrom != "" && listTo != "" {
		return dateutil.CustomRange(listFrom, listTo)
	}

	return dateutil.DateRange{}, fmt.Errorf("--from and --to must both be specified")
}

func processdays(client *slack.Client, days []time.Time, parallel int) []collector.DayResult {
	if parallel < 1 {
		parallel = 1
	}

	// Create work channel
	work := make(chan time.Time, len(days))
	for _, day := range days {
		work <- day
	}
	close(work)

	// Create results channel
	results := make(chan collector.DayResult, len(days))

	// Start workers
	var wg sync.WaitGroup
	for i := 0; i < parallel; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for day := range work {
				result := processDay(client, day)
				results <- result
			}
		}()
	}

	// Wait for all workers to complete
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect results
	var allResults []collector.DayResult
	for result := range results {
		allResults = append(allResults, result)
	}

	return allResults
}

func processDay(client *slack.Client, day time.Time) collector.DayResult {
	opts := collector.ListOptions{
		Date:            day,
		Author:          listAuthor,
		Mentions:        listMentions,
		Channels:        listChannels,
		ExcludeChannels: listExcludeChannels,
		WithThread:      listThread,
	}

	result, err := collector.List(client, opts)
	if err != nil {
		return collector.DayResult{
			Date:  day,
			Error: err,
		}
	}

	// Write to file
	outputPath := dateutil.OutputPath(day)
	writer, err := output.NewFileWriter(outputPath)
	if err != nil {
		return collector.DayResult{
			Date:  day,
			Error: fmt.Errorf("failed to create output file: %w", err),
		}
	}
	defer writer.Close()

	if len(result.Threads) > 0 {
		if err := writer.Write(result.Threads); err != nil {
			return collector.DayResult{
				Date:  day,
				Error: fmt.Errorf("failed to write output: %w", err),
			}
		}
	} else {
		// Write empty array
		if err := writer.Write([]interface{}{}); err != nil {
			return collector.DayResult{
				Date:  day,
				Error: fmt.Errorf("failed to write output: %w", err),
			}
		}
	}

	return *result
}
