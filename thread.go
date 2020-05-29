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
	params config.Config,
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
	if params.Parameters.MaxQueryLength > 0 && len(options.QueryString) > params.Parameters.MaxQueryLength {
		options.QueryString = options.QueryString[:params.Parameters.MaxQueryLength]
	}

	// Convert string to array of words. If labels are defined, extract those labels into different categories.
	parsed := methods.ParseQueryString(options.QueryString, params.Labels)

	selector, err := selectMethods.BuildSelectQuery(
		parsed,
		params.Bucket,
		options.Fields,
		options.Joins,
		params.LabelsOptions,
	)

	if err != nil {
		return "", err
	}

	where := whereMethods.BuildWhereQuery(parsed, options.Where, params.LabelsOptions)
	order := orderMethods.BuildOrderQuery(parsed, options.Order, params.LabelsOptions)

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

	if params.Parameters.Debug {
		fmt.Println(queryString)
	}

	results, clusterErr := params.Cluster.Query(queryString, nil)
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