package nefts

import (
	"fmt"
	"github.com/Kushuh/nefts-go/config"
	"github.com/Kushuh/nefts-go/methods"
	"github.com/Kushuh/nefts-go/methods/orderMethods"
	"github.com/Kushuh/nefts-go/methods/selectMethods"
	"github.com/Kushuh/nefts-go/methods/whereMethods"
	"strconv"
)

func BuildQueryString(
	start uint64,
	end uint64,
	options config.Options,
) (string, *config.Error) {
	// Allow negative start value but cap it to 0.
	if start < 0 {
		start = 0
	}

	if end < start {
		return "", &config.Error{
			Status: 400,
			Message: fmt.Sprintf("%q end point cannot be lower than %q start point !", end, start),
		}
	}

	// Negative value disable max query length. Default is 1000.
	if options.Config.Parameters.MaxQueryLength > 0 && len(options.QueryString) > options.Config.Parameters.MaxQueryLength {
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
		strconv.FormatUint(end - start, 10) + " OFFSET " + strconv.FormatUint(start, 10)

	return finalQuery, nil
}

func Thread(
	params config.Config,
	start uint64,
	end uint64,
	options config.Options,
) (*config.QueryResults, *config.Error) {
	queryString, err := BuildQueryString(params, start, end, options)

	if err != nil {
		return nil, err
	}

	if options.Config.Parameters.Debug {
		fmt.Println(queryString)
	}

	results, clusterErr := options.Config.Cluster.Query(queryString, nil)
	if clusterErr != nil {
		return nil, &config.Error{
			Status: 500,
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
	output.Boundaries.End = start + uint64(len(output.Results))

	output.Flags.BeginningOfResults = start == 0
	output.Flags.EndOfResults = uint64(len(output.Results)) < (end - start)

	return &output, nil
}