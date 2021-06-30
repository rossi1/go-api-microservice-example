package cmd

/*
import (
	"context"
	"fmt"

	"github.com/olivere/elastic"
	"github.com/spf13/cobra"
)

const mappings = `{
	"settings":{
	   "number_of_shards":1,
	   "number_of_replicas":0
	},
	"analysis":{
	   "analyzer":{
		  "standard_analyzer":{
			 "filter":[
				"lowercase",
				"snowball"
			 ],
			 "token_chars":[
				"letter",
				"digit"
			 ],
			 "type":"custom",
			 "tokenizer":"standard"
		  }
	   }
	},
	"mappings":{
	   "properties":{
		  "name":{
			 "type":"text",
			 "fields":{
				"suggest":{
				   "type":"completion",
				   "analyzer":"simple",
				   "preserve_separators":true,
				   "preserve_position_increments":true,
				   "max_input_length":50
				}
			 },
			 "analyzer":"standard_analyzer"
		  },
		  "product":{
			 "type":"nested",
			 "properties":{
				"name":{
				   "type":"text",
				   "fields":{
					  "suggest":{
						 "type":"completion",
						 "analyzer":"simple",
						 "preserve_separators":true,
						 "preserve_position_increments":true,
						 "max_input_length":50
					  }
				   },
				   "analyzer":"standard_analyzer"
				},
				"description":{
				   "type":"text",
				   "anayzer":"standard_analyzer"
				},
				"weight":{
				   "type":"text",
				   "anayzer":"standard_analyzer"
				},
				"tax":{
				   "type":"text",
				   "anayzer":"standard_analyzer"
				},
				"bar_code":{
				   "type":"text",
				   "anayzer":"standard_analyzer"
				},
				"discount":{
				   "type":"text",
				   "anayzer":"standard_analyzer"
				},
				"image":{
				   "type":"text",
				   "anayzer":"standard_analyzer"
				},
				"expires":{
				   "type":"date",
				   "null_value":"NULL",
				   "format":"strict_date_optional_time||epoch_second"
				}
			 }
		  }
	   }
	}
 }`

func Execute(ctx context.Context, es *elastic.Client) error {

	// CMD for the program
	var CMD = &cobra.Command{
		Use:     "search_index",
		Short:   "search_index CLI",
		Long:    "search_index CLI",
		Version: "v0.1.0",
	}
	createIndexCMD, err := createIndexCmd(es)
	if err != nil {
		return fmt.Errorf("failed to run search index: %w", err)
	}

	CMD.AddCommand(createIndexCMD)
	if err := CMD.ExecuteContext(ctx); err != nil {
		return err
	}
	return nil

}

func createIndexCmd(es *elastic.Client) (*cobra.Command, error) {
	var startCmd = &cobra.Command{
		Use:   "create",
		Short: "create index",
	}

	startCmd.RunE = func(cmd *cobra.Command, args []string) error {
		if err := createIndex(cmd.Context(), es); err != nil {
			return err
		}
		return nil
	}
	return startCmd, nil

}

func createIndex(ctx context.Context, client *elastic.Client) error {

	exists, err := client.IndexExists("inventories").Do(ctx)

	if err != nil {
	}

	if !exists {
		_, err := client.CreateIndex("inventories").BodyString(mappings).Do(ctx)
		if err != nil {
			return err
		}
	}
	return nil

}
*/
