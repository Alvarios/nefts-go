package nefts

import (
	"fmt"
	"github.com/Alvarios/nefts-go/config"
	"github.com/Alvarios/nefts-go/methods"
	"github.com/Alvarios/nefts-go/methods/orderMethods"
	"github.com/Alvarios/nefts-go/methods/selectMethods"
	"github.com/Alvarios/nefts-go/methods/whereMethods"
	"strconv"
)

func BuildQueryString(
	start int64,
	end int64,
	options config.Options,
) (string, *config.Error) {
	// Allow negative start value but cap it to 0.
	if start < 0 {
		start = 0
	}

	if end < start {
		return "", &config.Error{
			Code:    400,
			Message: fmt.Sprintf("%q end point cannot be lower than %q start point !", end, start),
		}
	}

	// Negative value disable max query length. Default is 1000.
	if len(options.QueryString) > options.Config.Parameters.MaxQueryLength {
		options.QueryString = options.QueryString[:options.Config.Parameters.MaxQueryLength]
	}

	// Convert string to array of words. If labels are defined, extract those labels into different categories.
	parsed := methods.ParseQueryString(options.QueryString, options.Labels)
	selector, err := selectMethods.BuildSelectQuery(
		parsed,
		options.Config.Bucket,
		options.Fields,
		options.Joins,
		options.LabelsOptions,
	)

	if err != nil {
		return "", err
	}

	where := whereMethods.BuildWhereQuery(parsed, options.Where, options.LabelsOptions)
	order := orderMethods.BuildOrderQuery(parsed, options.Order, options.LabelsOptions)

	finalQuery := selector + "\n\t" + where + "\n\t" + order + "\n\tLIMIT " +
		strconv.FormatInt(end - start, 10) + " OFFSET " + strconv.FormatInt(start, 10)

	return finalQuery, nil
}

func Thread(
	start int64,
	end int64,
	options config.Options,
) (*config.QueryResults, *config.Error) {
	queryString, err := BuildQueryString(start, end, options)

	if err != nil {
		return nil, err
	}

	if options.Config.Parameters.Debug {
		fmt.Println(queryString)
	}

	results, clusterErr := options.Config.Cluster.Query(queryString, nil)
	if clusterErr != nil {
		return nil, &config.Error{
			Code:    500,
			Message: clusterErr.Error(),
		}
	}

	var output config.QueryResults

	for results.Next() {
		var result interface{}
		err := results.Row(&result)
		if err != nil {
			panic(err)
		}

		output.Results = append(output.Results, result)
	}

	output.Boundaries.Start = start
	output.Boundaries.End = start + int64(len(output.Results))

	output.Flags.BeginningOfResults = start == 0
	output.Flags.EndOfResults = int64(len(output.Results)) < (end - start)

	return &output, nil
}