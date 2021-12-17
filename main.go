package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"gofind/find"
	"os"
	"time"
)

var (
	contains     string
	iContains    string
	sizeAtLeast  int
	sizeAtMost   int
	maxDepth     int
	extension    string
	startsWith   string
	endsWith     string
	regex        string
	notContains  string
	iNotContains string
)

func main() {

	var rootCmd = &cobra.Command{
		Use:   "gofind [flags] PATH",
		Short: "short",
		Long:  "gofind searches the file system for files that match your filters",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			run(args[0])
		},
	}

	rootCmd.Flags().StringVarP(&contains, "contains", "c", "", "string that the file name contains")
	rootCmd.Flags().StringVarP(&iContains, "icontains", "i", "", "case insensitive string that the file name contains")
	rootCmd.Flags().IntVar(&sizeAtLeast, "size-at-least", 0, "min file size")
	rootCmd.Flags().IntVar(&sizeAtMost, "size-at-most", 0, "max file size")
	rootCmd.Flags().IntVarP(&maxDepth, "max-depth", "d", 0, "max recursion depth to search")
	rootCmd.Flags().StringVarP(&extension, "ext", "e", "", "file extension")
	rootCmd.Flags().StringVar(&startsWith, "starts-with", "", "string that the file name starts with")
	rootCmd.Flags().StringVar(&endsWith, "ends-with", "", "string that the file name ends with")
	rootCmd.Flags().StringVarP(&regex, "regex", "r", "", "regex pattern that the file name matches")
	rootCmd.Flags().StringVar(&notContains, "not-contains", "", "string that the file name doesn't contain")
	rootCmd.Flags().StringVar(&iNotContains, "inot-contains", "", "case insensitive string that the file name doesn't contain")

	_ = rootCmd.Execute()

}

func run(rootPath string) {
	f := find.NewFinder(rootPath)
	if contains != "" {
		f.Contains(contains)
	}
	if iContains != "" {
		f.ContainsInsensitive(iContains)
	}
	if sizeAtLeast != 0 {
		f.SizeAtLeast(sizeAtLeast)
	}
	if sizeAtMost != 0 {
		f.SizeAtMost(sizeAtMost)
	}
	if maxDepth != 0 {
		f.MaxDepth(uint(maxDepth))
	}
	if extension != "" {
		f.EndsWith(extension)
	}
	if startsWith != "" {
		f.StartsWith(startsWith)
	}
	if endsWith != "" {
		f.EndsWith(endsWith)
	}
	if regex != "" {
		f.Regex(regex)
	}
	if notContains != "" {
		f.DoesntContain(notContains)
	}
	if iNotContains != "" {
		f.DoesntContainsInsensitive(iNotContains)
	}

	f.PrintMatches()
	start := time.Now()
	matches, err := f.Find()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
	}
	fmt.Printf("%d matches\n", len(matches))
	fmt.Printf("took %f seconds\n", time.Since(start).Seconds())
}
