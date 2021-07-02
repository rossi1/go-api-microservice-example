package main

import (
	"context"
	"fmt"
	"github/rossi1/go-api-microservice-example/cmd/internal"
	"log"
	"os"

	"github.com/olivere/elastic/v6"
	"github.com/spf13/cobra"
)

const mappings = `{
	"settings":{
	   "number_of_shards":1,
	   "number_of_replicas":0,
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
	 }
	},
	
	"mappings":{
	  "_doc": {
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
		       "name":{"type":"text",
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
				   "analyzer":"standard_analyzer"
				},
				"weight":{
				   "type":"text",
				   "analyzer":"standard_analyzer"
				},
				"tax":{
				   "type":"text",
				   "analyzer":"standard_analyzer"
				},
				"bar_code":{
				   "type":"text",
				   "analyzer":"standard_analyzer"
				},
				"discount":{
				   "type":"text",
				   "analyzer":"standard_analyzer"
				},
				"image":{
				   "type":"text",
				   "analyzer":"standard_analyzer"
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
 }
}`

func Execute(ctx context.Context, es *elastic.Client) (string, error) {

	// CMD for the program
	var CMD = &cobra.Command{
		Use:     "search_index",
		Short:   "search_index",
		Long:    "search_index CLI",
		Version: "v0.1.0",
	}
	createIndexCMD, err := createIndexCmd(es)

	if err != nil {
		return "", fmt.Errorf("failed to create search index: %w", err)
	}

	destroyIndexCMD, err := deleteIndexCmd(es)

	if err != nil {
		return "", fmt.Errorf("failed to destroy search index: %w", err)
	}

	verifyIndexCMD, err := indexExistsCmd(es)

	if err != nil {
		return "", err
	}

	CMD.AddCommand(createIndexCMD)
	CMD.AddCommand(destroyIndexCMD)
	CMD.AddCommand(verifyIndexCMD)
	if err := CMD.ExecuteContext(ctx); err != nil {
		return "", err
	}

	return "sucess", nil

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

func deleteIndexCmd(es *elastic.Client) (*cobra.Command, error) {
	var startCmd = &cobra.Command{
		Use:   "destroy",
		Short: "destroy index",
	}

	startCmd.RunE = func(cmd *cobra.Command, args []string) error {
		if err := deleteIndex(cmd.Context(), es); err != nil {
			return err
		}
		return nil
	}
	return startCmd, nil

}

func indexExistsCmd(es *elastic.Client) (*cobra.Command, error) {
	var startCmd = &cobra.Command{
		Use:   "check",
		Short: "check index",
	}

	startCmd.RunE = func(cmd *cobra.Command, args []string) error {
		if err := indexExists(cmd.Context(), es); err != nil {
			return err
		}
		return nil
	}
	return startCmd, nil

}

func createIndex(ctx context.Context, client *elastic.Client) error {

	exists, err := client.IndexExists("inventories").Do(ctx)

	if err != nil {
		return err
	}

	if !exists {
		_, err := client.CreateIndex("inventories").BodyString(mappings).Do(ctx)
		if err != nil {
			return err
		}
	}
	return nil

}

func deleteIndex(ctx context.Context, client *elastic.Client) error {

	_, err := client.DeleteIndex("inventories").Do(ctx)

	if err != nil {
		return err
	}
	return nil

}

func indexExists(ctx context.Context, client *elastic.Client) error {

	exists, err := client.IndexExists("inventories").Do(ctx)

	if err != nil {
		return err
	}

	if !exists {
		return fmt.Errorf("%s doesn't exist", "inventories")
	}

	return nil

}

func main() {
	ctx := context.Background()

	cfg, err := internal.GetConfig(ctx)

	if err != nil {
		log.Fatalf("config error: %s", err)
		os.Exit(1)
	}

	es, err := internal.NewElasticSearch(cfg)

	if err != nil {
		log.Fatalf("elasticsearch error: %s", err)
		os.Exit(1)
	}
	info, code, err := es.Ping("http://127.0.0.1:9200").Do(ctx)

	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Printf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)

	msg, err := Execute(ctx, es)
	if err != nil {
		log.Fatal(err.Error())
		os.Exit(1)

	}
	fmt.Println(msg)
}
